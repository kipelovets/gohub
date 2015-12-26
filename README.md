gohub
=====

See original README at https://github.com/marina-lab/gohub and https://github.com/adjust/gohub

## Differences in this fork

### Signature check

This fork adds additional security by providing a check for GitHub signature, provide GitHub webhook secret with --secret command line option.

### Hooks for different branches

There is also an ability to define different hooks for different branches of the same repository.

### Stdout log

You can provide "-" as logfile name to write logs to STDOUT

### Docker Compose test environment

You can quickly spawn a test environment with:

    docker-compose up -d
    
This command compiles sources from current dir and launches `gohub` binary with `example.json` as config file (see `Dockerfile` and `docker-compose.yml` for details)

If you use `dinghy` you can access `http://gohub_test.docker:7654` or launch `test_hook.sh`

## License

This Software is licensed under the MIT License.

Copyright (c) 2012 adeven GmbH, 
http://www.adeven.com

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
"Software"), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
