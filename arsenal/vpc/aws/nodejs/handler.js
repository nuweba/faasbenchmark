function getDuration() {
    var startTime = process.hrtime();
    var end = process.hrtime(startTime);
    return end[1] + (end[0] * 1e9);
}

function isWarm() {
    var is_warm = process.env.warm ? true : false;
    process.env.warm = true;
    return is_warm;
}

exports.handler = async () => {
    return {
        "reused": isWarm(),
        "duration": getDuration(),
    };
};
