const buf = Buffer.allocUnsafe(10 * 1024 * 1024);
const Storage = require('@google-cloud/storage');
const storage = Storage();

async function upload() {
    let bucket = storage.bucket(process.env.TEST_BUCKET);
    let file = bucket.file('faastest' + (new Date).getTime());
    await new Promise((resolve) => file.save(buf.toString('binary'), () => resolve()));
}

function isWarm() {
    var is_warm = process.env.warm ? true : false;
    process.env.warm = true;
    return is_warm;
}

function getDuration(startTime) {
    var end = process.hrtime(startTime);
    return end[1] + (end[0] * 1e9);
}

function getLevel(event) {
    return {}
}

function getParameters(event) {
    return getLevel(event);
}

async function runTest(intensityLevel) {
    await upload()
}

async function main(req, res) {
    let startTime = process.hrtime();
    let params = getParameters(req);
    if (params.error) {
        return {"error": params.error}
    }

    await runTest(params);

    var reused = isWarm();
    var duration = getDuration(startTime);

    res.status(200).send({
        reused: reused,
        duration: duration
    });

}

exports.test128 = main;
exports.test256 = main;
exports.test512 = main;
exports.test1024 = main;
exports.test2048 = main;