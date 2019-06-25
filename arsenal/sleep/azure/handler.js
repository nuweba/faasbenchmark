'use strict';

/* eslint-disable no-param-reassign */

var wait = ms => new Promise((r, j) => setTimeout(r, ms));

async function sleep(sleep_time) {
    let startTime = process.hrtime();
    await wait(sleep_time);
    let end = process.hrtime(startTime);
    return end[1] + (end[0] * 1e9);
}

function isWarm() {
    var is_warm = process.env.warm ? true : false;
    process.env.warm = true;
    return is_warm;
}


module.exports.hello = async function (context) {

    const sleep_time = context.req.query["sleep"] ? parseInt(context.req.query["sleep"]) : 200;

    return {
        "reused": isWarm(),
        "duration": await sleep(sleep_time)
    }
};
