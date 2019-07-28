'use strict';

var getters = ['siblings', 'index', 'first', 'last', 'prev', 'next'];
var utils = require('./utils');

/**
 * Create a new AST `Node` with the given `val` and `type`.
 *
 * ```js
 * var node = new Node('*', 'Star');
 * var node = new Node({type: 'star', val: '*'});
 * ```
 * @name Node
 * @param {String|Object} `val` Pass a matched substring, or an object to merge onto the node.
 * @param {String} `type` The node type to use when `val` is a string.
 * @return {Object} node instance
 * @api public
 */

var Node = exports = module.exports = function Node(val, type) {
  this.define('isNode', true);
  this.type = null;

  if (utils.isObject(val)) {
    for (var key in val) {
      if (getters.indexOf(key) === -1) {
        this[key] = val[key];
      }
    }
  } else {
    this.type = type;
    this.val = val;
  }
};

/**
 * Define a non-enumberable property on the node instance.
 *
 * ```js
 * var node = new Node();
 * node.define('foo', 'something non-enumerable');
 * ```
 * @name .define
 * @param {String} `name`
 * @param {any} `val`
 * @return {Object} returns the node instance
 * @api public
 */

Node.prototype.define = function(name, val) {
  utils.define(this, name, val);
  return this;
};

/**
 * Given node `foo` and node `bar`, push node `bar` onto `foo.nodes`, and
 * set `foo` as `bar.parent`.
 *
 * ```js
 * var foo = new Node({type: 'foo'});
 * var bar = new Node({type: 'bar'});
 * foo.pushNode(bar);
 * ```
 * @name .pushNode
 * @param {Object} `node`
 * @return {undefined}
 * @api public
 */

Node.prototype.pushNode = function(node) {
  this.nodes = this.nodes || [];
  utils.define(node, 'parent', this);
  this.nodes.push(node);
};

/**
 * Alias for [pushNode](#pushNode) for backwards compatibility with 0.1.0.
 * @name .addNode
 * @api public
 */

Node.prototype.addNode = function(node) {
  return this.pushNode(node);
};

/**
 * Given node `foo` and node `bar`, unshift node `bar` onto `foo.nodes`, and
 * set `foo` as `bar.parent`.
 *
 * ```js
 * var foo = new Node({type: 'foo'});
 * var bar = new Node({type: 'bar'});
 * foo.unshiftNode(bar);
 * ```
 * @name .unshiftNode
 * @param {Object} `node`
 * @return {undefined}
 * @api public
 */

Node.prototype.unshiftNode = function(node) {
  this.nodes = this.nodes || [];
  utils.define(node, 'parent', this);
  this.nodes.unshift(node);
};

/**
 * Get the first child node from `node.nodes` that matches the given `type`.
 * If `type` is a number, the child node at that index is returned.
 *
 * ```js
 * var child = node.getNode(1); //<= index of the node to get
 * var child = node.getNode('foo');
 * var child = node.getNode(/^(foo|bar)$/);
 * var child = node.getNode(['foo', 'bar']);
 * ```
 * @name .getNode
 * @param {String} `type`
 * @return {Object} Returns a child node or undefined.
 * @api public
 */

Node.prototype.getNode = function(type) {
  return utils.su.getNode(this.nodes, type);
};

/**
 * Return true if the node is the given `type`.
 *
 * ```js
 * var node = new Node({type: 'bar'});
 * cosole.log(node.isType('foo'));          // false
 * cosole.log(node.isType(/^(foo|bar)$/));  // true
 * cosole.log(node.isType(['foo', 'bar'])); // true
 * ```
 * @name .isType
 * @param {String} `type`
 * @return {Boolean}
 * @api public
 */

Node.prototype.isType = function(type) {
  return utils.su.isType(this, type);
};

/**
 * Return true if the `node.nodes` has the given `type`.
 *
 * ```js
 * var foo = new Node({type: 'foo'});
 * var bar = new Node({type: 'bar'});
 * foo.pushNode(bar);
 *
 * cosole.log(foo.hasType('qux'));          // false
 * cosole.log(foo.hasType(/^(qux|bar)$/));  // true
 * cosole.log(foo.hasType(['qux', 'bar'])); // true
 * ```
 * @name .hasType
 * @param {String} `type`
 * @return {Boolean}
 * @api public
 */

Node.prototype.hasType = function(type) {
  return utils.su.hasType(this, type);
};

/**
 * Get the siblings array, or `null` if it doesn't exist.
 *
 * ```js
 * var foo = new Node({type: 'foo'});
 * var bar = new Node({type: 'bar'});
 * var baz = new Node({type: 'baz'});
 * foo.pushNode(bar);
 * foo.pushNode(baz);
 *
 * console.log(bar.siblings.length) // 2
 * console.log(baz.siblings.length) // 2
 * ```
 * @name .siblings
 * @return {Array}
 * @api public
 */

Object.defineProperty(Node.prototype, 'siblings', {
  get: function() {
    return this.parent ? this.parent.nodes : null;
  }
});

/**
 * Get the previous node from the siblings array or `null`.
 *
 * ```js
 * var foo = new Node({type: 'foo'});
 * var bar = new Node({type: 'bar'});
 * var baz = new Node({type: 'baz'});
 * foo.pushNode(bar);
 * foo.pushNode(baz);
 *
 * console.log(baz.prev.type) // 'bar'
 * ```
 * @name .prev
 * @return {Object}
 * @api public
 */

Object.defineProperty(Node.prototype, 'prev', {
  get: function() {
    return this.parent && this.siblings
      ? this.siblings[this.index - 1] || this.parent.prev
      : null;
  }
});

/**
 * Get the siblings array, or `null` if it doesn't exist.
 *
 * ```js
 * var foo = new Node({type: 'foo'});
 * var bar = new Node({type: 'bar'});
 * var baz = new Node({type: 'baz'});
 * foo.pushNode(bar);
 * foo.pushNode(baz);
 *
 * console.log(bar.siblings.length) // 2
 * console.log(baz.siblings.length) // 2
 * ```
 * @name .next
 * @return {Object}
 * @api public
 */

Object.defineProperty(Node.prototype, 'next', {
  get: function() {
    return this.parent && this.siblings
      ? this.siblings[this.index + 1] || this.parent.next
      : null;
  }
});

/**
 * Get the node's current index from `node.parent.nodes`.
 * This should always be correct, even when the parent adds nodes.
 *
 * ```js
 * var foo = new Node({type: 'foo'});
 * var bar = new Node({type: 'bar'});
 * var baz = new Node({type: 'baz'});
 * var qux = new Node({type: 'qux'});
 * foo.pushNode(bar);
 * foo.pushNode(baz);
 * foo.unshiftNode(qux);
 *
 * console.log(bar.index) // 1
 * console.log(baz.index) // 2
 * console.log(qux.index) // 0
 * ```
 * @name .index
 * @return {Number}
 * @api public
 */

Object.defineProperty(Node.prototype, 'index', {
  get: function() {
    return this.siblings ? this.siblings.indexOf(this) : -1;
  }
});

/**
 * Get the first node from `node.nodes`.
 *
 * ```js
 * var foo = new Node({type: 'foo'});
 * var bar = new Node({type: 'bar'});
 * var baz = new Node({type: 'baz'});
 * var qux = new Node({type: 'qux'});
 * foo.pushNode(bar);
 * foo.pushNode(baz);
 * foo.pushNode(qux);
 *
 * console.log(foo.first.type) // 'bar'
 * ```
 * @name .first
 * @return {Object} The first node, or undefiend
 * @api public
 */

Object.defineProperty(Node.prototype, 'first', {
  get: function() {
    return utils.su.arrayify(this.nodes)[0];
  }
});

/**
 * Get the last node from `node.nodes`.
 *
 * ```js
 * var foo = new Node({type: 'foo'});
 * var bar = new Node({type: 'bar'});
 * var baz = new Node({type: 'baz'});
 * var qux = new Node({type: 'qux'});
 * foo.pushNode(bar);
 * foo.pushNode(baz);
 * foo.pushNode(qux);
 *
 * console.log(foo.last.type) // 'qux'
 * ```
 * @name .last
 * @return {Object} The last node, or undefiend
 * @api public
 */

Object.defineProperty(Node.prototype, 'last', {
  get: function() {
    return utils.su.last(utils.su.arrayify(this.nodes));
  }
});
