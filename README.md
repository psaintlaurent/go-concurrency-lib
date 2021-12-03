I'm re-writing some basic libraries that I use in other languages to ensure I have a simple set of libs for doing basic concurrency.  I don't want to hve to re-write this.  This includes single parallel process, multi parallel processing, pipelines, task scheduling and task status updates.

TODO: Add different controller types for handling job routes.
TODO: Handle job retries.
TODO: Allow for job function definitions (see job.fn).
TODO: Restructure the code to remove queuing from http routes.
TODO: Use Little's Law to help predict job completion time.
TODO: Keep track of job success, failure, completion, progress rates and times.
TODO: Allow for job versioning so that metrics are tied to a specific job version.
TODO: Allow for job templates to make it easier to define job progress checkpoints


