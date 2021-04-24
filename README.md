# PGM-index

A Go implementation of the PGM index

## What is a PGM index?

[Website](https://pgm.di.unipi.it/) | [Paper](http://www.vldb.org/pvldb/vol13/p1162-ferragina.pdf)

> The Piecewise Geometric Model index (PGM-index) is a data structure that enables fast lookup, predecessor, range searches and updates in arrays of billions of items using orders of magnitude less space than traditional indexes while providing the same worst-case query time guarantees.

Paolo Ferragina and Giorgio Vinciguerra. The PGM-index: a fully-dynamic compressed learned index with provable worst-case bounds. PVLDB, 13(8): 1162-1175, 2020.

**Credit**

- The Rotating Calipers algorithm is ported from https://github.com/bkiers/RotatingCalipers. Complete credit goes to its author.
- The Monotone Chain algorithm for the convex hull is from https://en.wikibooks.org/wiki/Algorithm_Implementation/Geometry/Convex_hull/Monotone_chain.