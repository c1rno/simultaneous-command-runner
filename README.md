# Simultaneous command runner

It's a simple tool which allow run simple shell commands
 like `helm install <smth>` and in parallel see logs
 of started items `kubectls logs -f -lrelease=<current-release>`,
 in CI of cource.

Main difference from `gnu parallel` or `xargs` is
 that there you have one "main" command and this
 solution stops while "main" command stops and it
 return error code if "main" command fails, ignore
 fails of "siblings".

About error codes there are one exceptions:
 - error codes may returns not from "main" command
   in "setup" stage, where no commands started

About messaging, it's solution redirects all command's
 output to stdout (stderr too) and write own messages
 to stderr only.

# Example

```
scr << EOF
timeout: 100
wait_all: false
run:
    cmd: "echo 🦁🐯"
    debug: true
watch:
    - cmd: "tail -F /dev/null"
      debug: false
      disable_time: true
      name: useless
    - cmd: "kubectl logs -f -lname=my-app"
      debug: true
      name: k8s
    - cmd: "ls -lah /var/log/"
      debug: true
      name: ls
EOF
```

Yep, there are few additional params - you may setup global timeout,
 or wait all siblings commands (in this case error codes still based
 on "main" cmd)

`debug` and `name` params it's about verbosity;
By default messages writes with time, you may disable it
 with specified param `disable_time: true`

# Docker how-to

I think, there is no need to use many images with env + `scr`,
 so I made one `scratch` images and you may reuse it as:

```
... you stuff here

FROM c1rno/simultaneous-command-runner as scr

FROM <your working image>

COPY --from=scr /usr/local/bin/scr /usr/local/bin/scr

... you stuff here

```

So, use it just as carrier (alpine users may feel little butthurt).
