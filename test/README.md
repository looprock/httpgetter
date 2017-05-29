# testing instructions

The simplest way to  set up a basic riemann server on localhost to test your configuration on a mac is:

## install the server

brew install riemann

brew install riemann-client

# run the server

in this directory, run:

riemann ./riemann.config

leave your terminal open so you can see the logs in stdout

# test the server

in another terminal window, something like:

riemann-client send -s critical -S foobar -h testhost

If everything is good, you should see something like the following in the terminal window you launched the server in:

INFO [2017-05-28 21:18:28,756] defaultEventExecutorGroup-2-4 - riemann.config - #riemann.codec.Event{:host testhost, :service foobar, :state critical, :description nil, :metric nil, :tags nil, :time 1.496024308756E9, :ttl nil, :x-client riemann-c-client}

If you do, congrats, you're running a riemann test instance on 127.0.0.1 port 5555
