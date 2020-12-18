'use strict';

// only for poc
// usage: node app.js some/path/to/config.yaml

const fs = require('fs');
const yaml = require('js-yaml');

// fabric
const fabricCaClient = require('fabric-ca-client');
const { Wallets } = require('fabric-network');

const FuncCallFailed = 0
const FuncCallOK = 1
const FuncCallIgnore = 2

// global variables
let MSPID = "";
let OrgName = "";
let adminUserId = "";

const express = require('express')
const app = express()
app.use(express.json())

const args = process.argv
if(args.length !== 3) {
    console.log('No config file');
    process.exit(1);
}

let config = parseConfig(args[2]);
console.log(`Host: ${config.express.host}`);
const walletPath = config.fabric.walletPath;
insurePathExist(walletPath);

// parse config file
function parseConfig(file) {
    try {
        let cfgContent = fs.readFileSync(file);
        let cfg = yaml.safeLoad(cfgContent, 'utf8');
        console.log(cfg);
        adminUserId = cfg.admin.username;
        return cfg;
    } catch (e) {
        console.log(`load config[${file}]failed, ${e}`);
        process.exit(1);
    }
}

function insurePathExist(path) {
    // check if path exists
    // if not, create it
    if(!fs.existsSync(path)) {
        // create it
        console.log(`Path: ${path} doesn't exist, create it...`);
        fs.mkdirSync(path);
    }
}

// create ca client
function getCaClient(cfg) {
    try {
        const ccpFile = cfg.fabric.ccPath;
        const caHost = cfg.fabric.caHost;
        let ccpContent = fs.readFileSync(ccpFile);
        let ccp = yaml.safeLoad(ccpContent, 'utf8');
        console.log(ccp);
        // some keys exist in connection.yaml
        const caInfo = ccp.certificateAuthorities[caHost];
        const caTLSCerts = caInfo.tlsCACerts.pem;
        const caClient = new fabricCaClient(caInfo.url, {trustedRoots: caTLSCerts, verify: false}, caInfo.caName);

        // write global values
        OrgName = ccp.client.organization;
        MSPID = ccp.organizations[OrgName].mspid;
        console.log(`OrgName: ${OrgName}, MSPID: ${MSPID}`);
        return caClient;
    } catch (e) {
        console.error(`create ca client failed, ${e}`);
        return null;
    }
}

// create file type wallet
async function getWallet(cfg) {
    try {
        const walletPath = cfg.fabric.walletPath;
        const wallet = await Wallets.newFileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);
        return wallet;
    } catch (e) {
        console.error(`get wallet failed, ${e}`);
        return null;
    }
}

async function registerAndEnrollUser(caClient, wallet, userId, passwd) {
    try {
        // check if exists
        const userIdentity = await wallet.get(userId);
        if(userIdentity) {
            console.log(`User ${userId} has exists in the wallet`);
            return FuncCallIgnore;
        }

        // get admin's identity
        const adminIdentity = await wallet.get(adminUserId);
        if(!adminIdentity) {
            console.error(`admin user not exist in the wallet, please enroll it first`);
            return FuncCallFailed;
        }

        const provider = wallet.getProviderRegistry().getProvider(adminIdentity.type);
        const adminUser = await provider.getUserContext(adminIdentity, adminUserId);
        const secret = await caClient.register({
            enrollmentID: userId,
            enrollmentSecret: passwd,
            role: 'client',
        }, adminUser);

        // enroll it
        await enrollUser(caClient, wallet, userId, passwd);
        console.log(`register and enroll user:${userId} success`);

    } catch (e) {
        console.error(`register and enroll user: ${userId} failed`);
    }
}

// enroll user
async function enrollUser(caClient, wallet, userId, passwd) {
    // first check if had already enrolled
    try {
        const identity = await wallet.get(userId);
        if(identity) {
            console.log(`An identity for user[${userId}] already exists in the wallet`);
            return;
        }

        const enrollment = await caClient.enroll({
            enrollmentID: userId,
            enrollmentSecret: passwd,
        });

        const x509Identity = {
            credentials: {
                certificate: enrollment.certificate,
                privateKey: enrollment.key.toBytes(),
            },
            mspId: MSPID,
            type: 'X.509',
        };

        await wallet.put(userId, x509Identity);
        console.log(`Successfully enrolled user ${userId} and imported it into the wallet`);
        return FuncCallOK;
    } catch (e) {
        console.error(`enroll user ${userId} failed, ${e}`);
        return FuncCallFailed;
    }
}

const hostname = config.express.hostname;
const port = config.express.port;

app.get('/', (req, res) => {
    res.json({"data": "pong"});
})

// register user
app.post('/api/ca/register', async (req, res) => {
    const caClient = getCaClient(config);
    const wallet = await getWallet(config);
    if(caClient == null || wallet == null) {
        console.error(`get ca client or wallet failed`);
        res.status(400).json({"msg": "get ca client or wallet failed"});
        return;
    }

    const reqJson = req.body;
    const userId = reqJson.userId || "";
    const passwd = reqJson.passwd || "";
    console.log(`userId: ${userId}, passwd: ${passwd}`);

    if(userId === "" || passwd === "") {
        res.status(400).json({"msg": "userId or passwd can't be null"});
        return;
    }

    const ret = await registerAndEnrollUser(caClient, wallet, userId, passwd);
    if(ret === FuncCallFailed) {
        res.status(400).json({"msg": "failed"});
    } else {
        res.json({"msg": "OK"});
    }
})

// enroll user
app.post('/api/ca/enroll', async (req, res) => {
    const caClient = getCaClient(config);
    const wallet = await getWallet(config);
    if(caClient == null || wallet == null) {
        console.error(`get ca client or wallet failed`);
        res.status(400).json({"msg": "get ca client or wallet failed"});
        return;
    }

    const reqJson = req.body;
    const userId = reqJson.userId || "";
    const passwd = reqJson.passwd || "";
    console.log(`userId: ${userId}, passwd: ${passwd}`);
    if(userId === "" || passwd === "") {
        res.status(400).json({"msg": "userId or passwd can't be null"});
        return;
    }

    const ret = await enrollUser(caClient, wallet, userId, passwd);
    if(ret === FuncCallFailed) {
        res.status(400).json({"msg": "failed"});
    } else {
        res.json({"msg": "OK"});
    }
})

app.listen(port, hostname, () => {
    console.log(`Listen at: http://localhost:${port}`)
})