## Overview

Vitess is a storage platform for scaling MySQL. It is optimized to run as
effectively in cloud architectures as it does on dedicated hardware. It combines
many important features of MySQL with the scalability of a NoSQL database.

> **Scalability**
>
> * Eliminates high memory overhead of MySQL connections
> * Lets your database grow without adding sharding logic to your application
> * Provides built-in sharding
> * Supports live resharding with minimal read-only downtime

<!-- -->

> **Performance**
>
> * Automatically rewrites queries that hurt database performance
> * Uses caching mechanism on proxy server to mediate queries and prevent
>   duplicate queries from simultaneously reaching your database

<!-- -->

> **Manageability**
>
> * Uses a lock server like ZooKeeper or etcd to track and administer servers,
>   letting your application be blissfully ignorant of database topology
> * Automatically handles functions like master failover and backups, minimizing
>   any necessary downtime

## History

Vitess has been a fundamental component of YouTube infrastructure since 2011.
This section briefly summarizes the sequence of events that led to Vitess'
creation:

1.  YouTube's MySQL database reached a point when peak traffic would soon
    exceed the database's serving capacity. To temporarily alleviate the
    problem, YouTube created a master database for write traffic and a
    replica database for read traffic.
1.  With demand for cat videos at an all-time high, read-only traffic was
    still high enough to overload the replica database. So YouTube added
    more replicas, again providing a temporary solution.
1.  Eventually, write traffic became too high for the master database to
    handle, requiring YouTube to shard data to handle incoming traffic.
    (Sharding would have also become necessary if the overall size of the
    database became too large for a single MySQL instance.)
1.  YouTube's application layer was modified so that before executing any
    database operation, the code could identify the right database shard
    to receive that particular query.

<p style="text-align: left;">
Vitess let YouTube remove that logic from the source code, introducing
a proxy between the application and the database to route and manage
database interactions. Since then, YouTube has scaled its user base
by a factor of more than 50, greatly increasing its capacity to serve
pages, process newly uploaded videos, and more. Even more importantly,
Vitess is a platform that continues to scale.
</p>

## Features

<p style="text-align: left;">Let's quickly summarize a few key Vitess features:
</p>

* **Connection pooling**<br>
    Each MySQL connection has a memory overhead, which is around 256KB
    in the default MySQL configuration in addition to a significant CPU
    overhead associated with obtaining the connection. Vitess' BSON-based
    protocol creates very lightweight connections (around 32KB per
    connection), enabling Vitess servers to easily handle thousands of
    connections. Vitess uses Go's awesome concurrency support to map
    these connections to a pool of MySQL connections.<br><br>

* **Shard management**<br>
    As your database grows, you will likely want to implement
    horizontal sharding. But MySQL doesn't natively support sharding,
    so you will need to write the sharding code yourself and add
    sharding logic to your app.<br><br>
    Vitess enables sharding with minimal read-only downtime. For
    example, it supports split replication, dividing the replication
    stream so that a future shard master only gets statements that 
    could affect rows in its new shard. Vitess can also accommodate
    a custom sharding scheme that you already have in place.
    <br><br>

* **Workflow management**<br>
    Vitess helps you manage the lifecycle of your database instances
    by automatically handling various scenarios like master failover
    and replication. Vitess also keeps track of all of the metadata
    about the cluster configuration so that the cluster view is always
    up-to-date and consistent for different clients. What's more,
    your app doesn't need to keep track of the database topology
    as it changes.

## Contact

Ask questions in the
[vitess@googlegroups.com](https://groups.google.com/forum/#!forum/vitess)
discussion forum.

Subscribe to
[vitess-announce@googlegroups.com](https://groups.google.com/forum/#!forum/vitess-announce)
for low-frequency updates like new features and releases.
