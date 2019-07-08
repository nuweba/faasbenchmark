const http = require('http');
const fs = require('fs');

async function networkIntensive(level) {
    var startTime = process.hrtime();

    var fileSizeInMB = Math.pow(10, level - 1);
    const writable = fs.createWriteStream('/dev/null');
    var res = await new Promise((resolve, reject) => http.get({
        host: `www.ovh.net`,
        port: 80,
        path: `/files/${fileSizeInMB}Mb.dat`
    }, (res) => {
        var download = res.pipe(writable);
        download.on('close', () => resolve(res));
    }));
    var end = process.hrtime(startTime);
    return end[1] + (end[0] * 1e9);
}

function isWarm() {
    var is_warm = process.env.warm ? true : false;
    process.env.warm = true;
    return is_warm;
}

exports.handler = async (event) => {
    let intensityLevel = event.level ? parseInt(event.level) : null;
    if (!intensityLevel || intensityLevel < 1) {
        return {"error": "invalid level parameter"}
    }

    return {
        "reused": isWarm(),
        "duration": await networkIntensive(intensityLevel)
    };
};
