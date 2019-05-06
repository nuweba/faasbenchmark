var wait = ms => new Promise((r, j)=>setTimeout(r, ms));

exports.hello = async (event) => {
    let startTime = process.hrtime();
	const sleep_time = event.sleep ? parseInt(event.sleep) : 200;
	await wait(sleep_time);
	let is_warm = process.env.warm ? true : false;
	process.env.warm = true;
	let end = process.hrtime(startTime);
	return {
	    "reused" : is_warm,
	    "duration" : end[1] + (end[0] * 1e9),
	};
};
