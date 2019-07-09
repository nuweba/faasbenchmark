const http = require('http');
const fs = require('fs');

async function networkIntensive(level) {
    var fileSizeInMB = Math.pow(10, level - 1);
    const writable = fs.createWriteStream('/dev/null');
    await new Promise((resolve) => http.get({
        host: `www.ovh.net`,
        port: 80,
        path: `/files/${fileSizeInMB}Mb.dat`
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

exports.handler = async (event) => {
    var startTime = process.hrtime();
    let intensityLevel = event.level ? parseInt(event.level) : null;
    if (!intensityLevel || intensityLevel < 1) {
        return {"error": "invalid level parameter"}
    }

    await networkIntensive(intensityLevel);

    let retval = {
        "reused": isWarm(),
    };

    var end = process.hrtime(startTime);
    retval.duration = end[1] + (end[0] * 1e9);
    return retval;
};
