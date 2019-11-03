args = argv();
iterations = str2num(args{1});
script_dir = fileparts(mfilename('fullpath'));
source(strcat(script_dir, "/", "iterate.m"));
iterate(iterations)
