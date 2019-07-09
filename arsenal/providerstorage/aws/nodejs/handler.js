var buf = Buffer.allocUnsafe(10 * 1024 * 1024);
const AWS = require('aws-sdk');
const s3 = new AWS.S3();

async function upload() {
    await s3.putObject({
        Bucket: process.env.TEST_BUCKET,
        Key: "faastest.dat",
        Body: buf,
    }).promise();
}

function isWarm() {
    var is_warm = process.env.warm ? true : false;
    process.env.warm = true;
    return is_warm;
}

exports.handler = async () => {
    var startTime = process.hrtime();

    await upload();
    let retval = {
        "reused": isWarm(),
    };

    var end = process.hrtime(startTime);
    retval.duration = end[1] + (end[0] * 1e9);
    return retval;
};

