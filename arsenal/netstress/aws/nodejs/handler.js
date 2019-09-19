const http = require('http');
const fs = require('fs');

const files = {1: '/files/1Mb.dat', 2: '/files/10Mb.dat', 3: '/files/100Mb.dat'};

async function networkIntensive(level) {
    const writable = fs.createWriteStream('/dev/null');
    await new Promise((resolve) => http.get({
        host: `www.ovh.net`,
        port: 80,
        path: files[level]
    }, (res) => {
        var download = res.pipe(writable);
        download.on('close', () => resolve(res));
    }));
}

function isWarm() {
    var is_warm = process.env.warm ? true : false;
    process.env.warm = true;
    return is_warm;
}

function getLevel(event) {
    let intensityLevel = event.level ? parseInt(event.level) : null;
    if (!intensityLevel || intensityLevel < 1) {
        return {"error": "invalid level parameter"};
    }
    return intensityLevel;
}

function getParameters(event) {
    return getLevel(event);
}

function getDuration(startTime) {
    var end = process.hrtime(startTime);
    return end[1] + (end[0] * 1e9);
}

async function runTest(intensityLevel){
    await networkIntensive(intensityLevel)
}

exports.handler = async (event) => {
    var startTime = process.hrtime();
    let params = getParameters(event);
    if (params.error) {
        return {"error": params.error}
    }

    await runTest(params);

    var reused = isWarm();
    var duration = getDuration(startTime);

    return {
        "reused": reused,
        "duration": duration
    };
};
