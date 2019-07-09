function isWarm() {
    var is_warm = process.env.warm ? true : false;
    process.env.warm = true;
    return is_warm;
}

exports.handler = async () => {
    var startTime = process.hrtime();

    let retval = {
        "reused": isWarm(),
    };

    var end = process.hrtime(startTime);
    retval.duration = end[1] + (end[0] * 1e9);
    return retval;
};
