var wait = ms => new Promise((r, j)=>setTimeout(r, ms));

async function testfunc(req, res) {
	let startTime = process.hrtime();
	const sleep_time = req.query.sleep ? parseInt(req.query.sleep) : null;
	await wait(sleep_time);
	let is_warm = process.env.warm ? true : false;
  	process.env.warm = true
	let end = process.hrtime(startTime);
	
	res.status(200).send({reused : is_warm, duration : end[1] + (end[0] * 1e9)});
	
}

exports.test128 = testfunc;
exports.test256 = testfunc;
exports.test512 = testfunc;
exports.test1024 = testfunc;
exports.test2048 = testfunc;