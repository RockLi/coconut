#coconut

Coconut is a foundation library fully written in Go which includes frequent used data structures and algorithms.


##Aims of this project

* Avoid to reinvent frequent used data structures and algorithms every time in different projects
* Avoid to import from different places for different data structure and algorithms
* Trying to be a good candidate as the foundation library


##Status of coconut

Before I release coconut v0.1, **donot** use this library in your project, 'cause: 

* I haven't add enough data structure and algorithms in
* Probably I will modify the API frequently
* Probably I will regroup the packages like merging different individual small packages into a more meaningful one
* Should have many bugs and places can be optimized


##Planned data strucutre and algorithms

_* means already done, but probably I will still refactor it. _

* Cache
  * LRU (*)
  * LFU (*)

* Scheduling 
  * RoundRobin (*)

* Bloom Filter
  * Standard Bloom Filter
  * Counting Bloom Filter
  * Scalable Bloom Filter

* Tree
  * B(+/*) Tree
  * AVL
  * Red Black Tree
  * Trie
  * R-trees

* Binary Search
  * Binary Search

* Hash
  * murmur3

* Set
  * Set

* Bitmap
  * Bitmap

* SkipList
  * Skip List

* FSM
  * FSM

* Utils
  * Gcd (*)

* _To be added_

##Docs

_Currently I haven't finished any docs even valuable comments inside the codes, I will find time to finish this._




