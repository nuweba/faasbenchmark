var wait = ms => new Promise((r, j) => setTimeout(r, ms));

async function sleep(sleep_time) {
    let startTime = process.hrtime();
    await wait(sleep_time);
    let end = process.hrtime(startTime);
    return end[1] + (end[0] * 1e9)
}

exports.hello = async (event) => {
    var is_warm = process.env.warm ? true : false;
    process.env.warm = true;
    const sleep_time = event.sleep ? parseInt(event.sleep) : 200;

    return {
        "reused": is_warm,
        "duration": sleep(sleep_time)
    };
};
