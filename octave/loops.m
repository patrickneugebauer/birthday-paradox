args = argv();
iterations = str2num(args{1});
script_dir = fileparts(mfilename('fullpath'));

% no need to source this, octave will read it as a function file not a script file
% https://octave.org/doc/v4.0.3/Function-Files.html#Function-Files
% https://octave.org/doc/v4.0.3/Script-Files.html#Script-Files
% source(strcat(script_dir, "/", "iterate.m"));

iterate(iterations)
