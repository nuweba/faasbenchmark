'use strict';

var utils = require('lazy-cache')(require);
var fn = require;
require = utils;

/**
 * Lazily required module dependencies
 */

require('define-property', 'define');
require('kind-of', 'typeOf');
require('snapdragon-node', 'Node');
require = fn;

/**
 * Expose `utils` modules
 */

module.exports = utils;
