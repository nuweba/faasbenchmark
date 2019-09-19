# snapdragon-util [![NPM version](https://img.shields.io/npm/v/snapdragon-util.svg?style=flat)](https://www.npmjs.com/package/snapdragon-util) [![NPM monthly downloads](https://img.shields.io/npm/dm/snapdragon-util.svg?style=flat)](https://npmjs.org/package/snapdragon-util)  [![NPM total downloads](https://img.shields.io/npm/dt/snapdragon-util.svg?style=flat)](https://npmjs.org/package/snapdragon-util) [![Linux Build Status](https://img.shields.io/travis/jonschlinkert/snapdragon-util.svg?style=flat&label=Travis)](https://travis-ci.org/jonschlinkert/snapdragon-util)

> Utilities for the snapdragon parser/compiler.

<details>
<summary><strong>Table of Contents</strong></summary>
- [Install](#install)
- [Usage](#usage)
- [API](#api)
- [About](#about)
</details>

## Install

Install with [npm](https://www.npmjs.com/):

```sh
$ npm install --save snapdragon-util
```

## Usage

```js
var util = require('snapdragon-util');
```

## API

### [.noop](index.js#L17)

Emit an empty string to effectively "skip" the string for the given `node`, but still emit the position and node type.

**Params**

* **{Object}**: node

**Example**

```js
// do nothing for beginning-of-string
snapdragon.compiler.set('bos', utils.noop);
```

### [.emit](index.js#L37)

Emit an empty string to effectively "skip" the string for the given `node`, but still emit the position and node type.

**Params**

* **{Object}**: node

**Example**

```js
snapdragon.compiler
  .set('i', function(node) {
    this.mapVisit(node);
  })
  .set('i.open', utils.emit('<i>'))
  .set('i.close', utils.emit('</i>'))
```

### [.toNoop](index.js#L54)

Converts an AST node into an empty `text` node and delete `node.nodes`.

**Params**

* `node` **{Object}**

**Example**

```js
utils.toNoop(node);
utils.toNoop(node, true); // convert `node.nodes` to an empty array
```

### [.visit](index.js#L84)

Visit `node` with the given `fn`. The built-in `.visit` method in snapdragon automatically calls registered compilers, this allows you to pass a visitor function.

**Params**

* `node` **{Object}**
* `fn` **{Function}**
* `returns` **{Object}**: returns the node

**Example**

```js
snapdragon.compiler
  .set('i', function(node) {
    exports.visit(node, function(node2) {
      // do stuff with "node2"
      return node2;
    });
  })
```

### [.mapVisit](index.js#L120)

Map [visit](#visit) with the given `fn` over an array of AST `nodes`.

**Example**

```js
snapdragon.compiler
  .set('i', function(node) {
    exports.mapVisit(node, function(node2) {
      // do stuff with "node2"
      return node2;
    });
  })
```

### [.wrapNodes](index.js#L152)

Wrap the given `node` with `*.open` and `*.close` tags.

**Params**

* `node` **{Object}**
* `filter` **{Function}**: Optionaly specify a filter function to exclude the node.
* `returns` **{undefined}**

### [.addOpen](index.js#L166)

Unshift an `*.open` node onto `node.nodes`.

**Params**

* `node` **{Object}**
* `filter` **{Function}**: Optionaly specify a filter function to exclude the node.
* `returns` **{undefined}**

### [.addClose](index.js#L185)

Push a `*.close` node onto `node.nodes`.

**Params**

* `node` **{Object}**
* `filter` **{Function}**: Optionaly specify a filter function to exclude the node.
* `returns` **{undefined}**

### [.pushNode](index.js#L210)

Push `node` onto `parent.nodes`.

**Params**

* `node` **{Object}**
* `filter` **{Function}**: Optionaly specify a filter function to exclude the node.
* `returns` **{undefined}**

**Example**

```js
var parent = new Node({type: 'foo'});
var node = new Node({type: 'bar'});
utils.pushNode(parent, node);
console.log(parent.nodes[0].type) // 'bar'
```

### [.unshiftNode](index.js#L230)

Unshift `node` onto `parent.nodes`.

**Params**

* `node` **{Object}**
* `returns` **{undefined}**

**Example**

```js
var parent = new Node({type: 'foo'});
var node = new Node({type: 'bar'});
utils.unshiftNode(parent, node);
console.log(parent.nodes[0].type) // 'bar'
```

### [.firstOfType](index.js#L306)

Return the first node from `nodes` of the given `type`

**Params**

* `nodes` **{Array}**
* `type` **{String}**
* `returns` **{Object}**: Returns a node, if found

**Example**

```js
snapdragon.set('div', function(node) {
 var textNode = exports.firstOfType(node.nodes, 'text');
 if (textNode) {
   // do stuff with text node
 }
});
```

### [.arrayify](index.js#L505)

Cast the given `val` to an array.

**Params**

* `val` **{any}**
* `returns` **{Array}**

### [.stringify](index.js#L518)

Convert the given `val` to a string by joining with `,`. Useful
for creating a selector from a list of strings.

**Params**

* `val` **{any}**
* `returns` **{Array}**

## About

### Contributing

Pull requests and stars are always welcome. For bugs and feature requests, [please create an issue](../../issues/new).

Please read the [contributing guide](.github/contributing.md) for advice on opening issues, pull requests, and coding standards.

### Building docs

_(This project's readme.md is generated by [verb](https://github.com/verbose/verb-generate-readme), please don't edit the readme directly. Any changes to the readme must be made in the [.verb.md](.verb.md) readme template.)_

To generate the readme, run the following command:

```sh
$ npm install -g verbose/verb#dev verb-generate-readme && verb
```

### Running tests

Running and reviewing unit tests is a great way to get familiarized with a library and its API. You can install dependencies and run tests with the following command:

```sh
$ npm install && npm test
```

### Author

**Jon Schlinkert**

* [github/jonschlinkert](https://github.com/jonschlinkert)
* [twitter/jonschlinkert](https://twitter.com/jonschlinkert)

### License

Copyright © 2017, [Jon Schlinkert](https://github.com/jonschlinkert).
MIT

***

_This file was generated by [verb-generate-readme](https://github.com/verbose/verb-generate-readme), v0.4.2, on February 15, 2017._