const fs = require('fs');

const DEFAULT_INTENSITY_LEVEL = 5;
const OUTPUT_FILE_PATH = "/tmp/current_value";

function ioIntensiveCalculation(baseNumber) {
    var startTime = process.hrtime();
    var iterationCount = 30 * Math.pow(baseNumber, 3) + 1000;
    var finalBufferLength = 0;
    var totalBytesRead = 0;
    for (var i = iterationCount; i >= 0; i--) {
        var bufToWrite = "Iteration Number #" + i + "wrote this message to this file, ";
        fs.appendFileSync(OUTPUT_FILE_PATH, bufToWrite);
        var newerFileBuf = fs.readFileSync(OUTPUT_FILE_PATH);
        finalBufferLength = newerFileBuf.length;
        totalBytesRead += newerFileBuf.length;
    }
    var end = process.hrtime(startTime);
    return end[1] + (end[0] * 1e9);
}

function isWarm() {
    var is_warm = process.env.warm ? true : false;
    process.env.warm = true;
    return is_warm;
}

exports.handler = async (event) => {
    let got_custom_level = event["level"] && event["level"] !== "0";
    let intensityLevel = got_custom_level ? parseInt(event["level"]) : DEFAULT_INTENSITY_LEVEL;

    return {
        "reused": isWarm(),
        "duration": ioIntensiveCalculation(intensityLevel)
    };
};

