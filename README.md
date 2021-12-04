I'm re-writing some basic libraries that I use in other languages to ensure I have a simple set of libs for doing basic concurrency.  I don't want to re-write this repeatedly.  This includes single parallel process, multi parallel processing, pipelines, task scheduling and task status updates.

TODO: Add different controller types for handling job routing and distribution.

TODO: Handle job retries.

TODO: Allow for job function definitions (see Job.fn).

TODO: split Type definitions and pointer receivers into separate files from main.go

TODO: Add GRPC

TODO: Restructure the code to remove queuing from http routes.

TODO: Use Little's Law to help predict job completion time.

TODO: Keep track of job success, failure, completion, progress rates and times.

TODO: Allow for job versioning so that metrics are tied to a specific job version.

TODO: Allow for job templates to make it easier to define job progress checkpoints


