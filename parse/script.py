import mdl
import sys
from os import remove, execlp, fork
def first_pass(commands):
    basename = "image"
    frames = 1
    has_vary = False
    for command in commands:
        if command["op"] == "frames":
            frames = int(command["args"][0])
        if command["op"] == "basename":
            basename = command["args"][0]
        if command["op"] == "vary":
            has_vary = True
    if frames == 1 and has_vary:
        print "VARY COMMAND WITH 0 FRAMES BOSS"
        sys.exit()
    return (frames, basename)
## returns an array of dictionaries with every knob value
def linear_func(x):
    return x
def second_pass(commands, num_frames, func):
    knobs = []
    for _ in range(num_frames):
        knobs.append({})
    # go through command, find a vary, and do the ting
    for command in commands:
        if command["op"] == "vary":
            knob = command["knob"]
            start_frame = command["args"][0]
            end_frame = command["args"][1]
            start_val = command["args"][2]
            end_val = command["args"][3]
            t_total = (end_frame-start_frame)
            delta_val = end_val-start_val
            for frame_num in range(int(num_frames)):
                if frame_num < start_frame or frame_num > end_frame:
                    knobs[frame_num][knob] = 1.0
                else:
                    t = (float(frame_num) - start_frame) / t_total
                    val =  start_val + delta_val * t
                    knobs[frame_num][knob] = val
    return knobs


def run(filename):
    """
    This function runs an mdl script
    """
    p = mdl.parseFile(filename)

    if p:
        (commands, symbols) = p
        (frames, basename) = first_pass(commands)
        knobs = second_pass(commands, frames, linear_func)
        print commands
        print knobs
        # we have knobs for each frame, this is where the fun starts
        # commands that can have knobs are move,rotate,scale
        knobbed_cmds = ["move","scale","rotate"]
        no_send = ["vary", "set", "frames", "basename"]
        DIRNAME = "img/"
        ctr = 0
        for frame in knobs:
            print "clear"
            for command in commands:
                if command["op"] in knobbed_cmds:
                    print (command["op"])
                    output_line = []
                    knob_mod = 1.0
                    if command["knob"] is not None:
                        knob_mod = frame[command["knob"]]
                    for arg in command["args"]:
                        if isinstance(arg, float):
                            arg *= knob_mod
                        output_line.append("%s " % str(arg))
                    print ''.join(output_line)
                 # non knobbed, just print it regularly with args on next line
                 # do not print commands that have no meaning to go
                elif command["op"] not in no_send:
                    print (command["op"])
                    output_line = []
                    if command["args"] is not None:
                        for arg in command["args"]:
                            output_line.append("%s " % str(arg))
                    print ''.join(output_line)
            ### SAVE FILE
            if frames > 1:
                print "save"
                print DIRNAME + "%s%03d.png" % (basename,ctr)
                ctr+=1
        ## save everything
        if frames > 1:
            name_arg = DIRNAME + basename + "*"
            name = basename+ ".gif"
            print "animate"
            print  name_arg + " " + name
    else:
        print "Parsing failed."
        return
