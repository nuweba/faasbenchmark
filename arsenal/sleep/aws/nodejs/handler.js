var wait = ms => new Promise((r, j) => setTimeout(r, ms));

function isWarm() {
    var is_warm = process.env.warm ? true : false;
    process.env.warm = true;
    return is_warm;
}

exports.hello = async (event) => {
    let startTime = process.hrtime();

    const sleep_time = event.sleep ? parseInt(event.sleep) : null;
	if (!sleep_time && sleep_time !== 0) {
		return {"error": "invalid sleep parameter"}
	}

    await wait(sleep_time);

    let retval = {
        "reused": isWarm(),
    };

    let end = process.hrtime(startTime);
    retval.duration = end[1] + (end[0] * 1e9);
    return retval;
};
