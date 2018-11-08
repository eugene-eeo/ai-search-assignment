import time

def charnum(astring):

    m = len(astring)
    dumlist = []
    if m != 0:
        for i in range(0,m):
            if ord(astring[i]) >= 48 and ord(astring[i]) <= 57:
                dumlist.append(astring[i])
        astring = ""
        if len(dumlist) != 0:
            for i in range(0,len(dumlist)):
                astring = astring + dumlist[i]
    if astring == "":
        astring = "0"
    return astring

def tc(inputfile,inputtourfile):

    f = open(inputfile,'r')  # read in data character by character
    x = f.read(1)            # and strip away rubbish
    flag = "good"
    d = ""
    while x != "":
        if ord(x) >= 44 and ord(x) <= 122:
            d = d + x
        x = f.read(1)
    f.close()

    start_of_name = d.find("NAME")  # check that NAME=<city file-name>, appears
    if start_of_name == -1:  
        flag = "error: no city-file name"
        trace.write(flag + "\n")
    else:
        comma_after_name = d.find(",",start_of_name)
        if comma_after_name == -1:
            flag = "error: no comma after NAME="
            trace.write(flag + "\n")
        else:
            name = d[start_of_name + 5:comma_after_name]  # get city name

    if flag[0] != "e":
        start_of_size = d.find("SIZE=",comma_after_name)  # check that SIZE=<number of cities>, appears
        if start_of_size == -1:
            flag = "error: no SIZE="
            trace.write(flag + "\n")
        else:
            comma_after_size = d.find(",",start_of_size)
            if comma_after_size == -1:
                flag = "error: no comma after SIZE="
                trace.write(flag+"\n")
            else:
                w = charnum(d[start_of_size + 5:comma_after_size])  # get number of sities
                n = int(w)
            
    if flag[0] != "e":
        start = comma_after_size + 1  # compute distances
        end = len(d)
        first_character = start
        row = [0]
        matrix = []
        current_row = 1
        current_column = 2
        if first_character >= end:
            jumpout = "yes"
        else:
            jumpout = "no"
        while jumpout == "no" and current_row <= n - 1:  # iteratively get numbers between commas
            next_comma = d.find(",",first_character)
            if next_comma != -1:
                step = int(charnum(d[first_character:next_comma]))
                first_character = next_comma + 1
            else:
                if first_character < end:
                    step = int(charnum(d[first_character:end]))
                    jumpout == "nearly"
                else:
                    jumpout = "yes"
            if jumpout != "yes":
                row.append(step)  # add distance to matrix
                current_column = current_column + 1
                if current_column == n + 1:    # add in dummy 0s
                    matrix.append(row)
                    row = [0]
                    current_row = current_row + 1
                    for k in range(1,current_row):
                        row.append(0)
                    current_column = current_row + 1
        if current_row <= n - 1:
            jumpout == "yes"
        if jumpout == "yes":
            flag = "error: too few distances supplied (up to [" + str(row) + "," + str(column) + "]"
            trace.write(flag + "\n")
        else:
            row = [0]
            for j in range(1,n):  # add on row of 0's
                row.append(0)
            matrix.append(row)
            for i in range(2,n + 1):  # make matrix symmetric
                for j in range(1,i):
                    matrix[i - 1][j - 1] = matrix[j - 1][i - 1]

    if flag[0] != "e":
        f = open(inputtourfile,'r')
        checklength = 0
        length = 0
        x = f.read(1)
        flag = "good"
        d = ""
        while x != "":
            if ord(x) >= 44 and ord(x) <= 122:
                    d = d + x
            x = f.read(1)
        f.close()

        start_of_name = d.find("NAME")  # look for city name
        if d.find("NAME=") == -1:
            flag = "error: no NAME= in tour-file"
            trace.write(flag + "\n")
        else:
            comma_after_name = d.find(",",d.find("NAME="))
            if comma_after_name == -1:
                flag = "error: no comma after NAME= in tour-file"
                trace.write(flag + "\n")
            else:
                tname = d[start_of_name + 5:comma_after_name]
                if name != tname:
                    flag = "error: name of city-file and tour-file don't match (city: " + name + " and tour: " + tname + ")"
                    trace.write(flag + "\n")
                    
        if flag[0] != "e":  # look for number of cities
            start_of_size = d.find("SIZE=",comma_after_name)
            if start_of_size == -1:
                flag = "error: no SIZE= in tour-file"
                trace.write(flag + "\n")
            else:
                comma_after_size = d.find(",",d.find("SIZE="))
                if comma_after_size == -1:
                    flag = "error: no comma after SIZE= in tour-file"
                    trace.write(flag + "\n")
                else:
                    tn = int(charnum(d[start_of_size + 5:comma_after_size]))
                    if tn != n:
                        flag = "error: numbers of cities mismatch (city size: " + str(n) + ", tour size: " + str(tn) + ")"
                        trace.write(flag + "\n")

        if flag[0] != "e":  # look for tour length
            start_of_length = d.find("LENGTH=")
            if start_of_length == -1:
                flag = "error: no LENGTH= in tour-file"
                trace.write(flag + "\n")
            else:
                comma_after_length = d.find(",",d.find("LENGTH="))
                if comma_after_length == -1:
                    flag = "error: no comma after LENGTH= in tour-file"
                    trace.write(flag + "\n")
                else:
                    length = int(charnum(d[start_of_length + 7:comma_after_length]))
                    
        if flag[0] != "e":  # check tour length
            start_digit = d.find(",",d.find("LENGTH=")) + 1
            end = len(d)
            city_count = 0
            tour = []
            if start_digit >= end:
                jumpout = "yes"
            else:
                jumpout = "no"
            while jumpout == "no" and city_count < n:  # iteratively get the next tour city
                next_comma = d.find(",",start_digit)
                if next_comma != -1:
                    step = int(charnum(d[start_digit:next_comma]))
                    if (step < 1) or (step > n):
                        jumpout = "yes"
                    else:
                        start_digit = next_comma + 1
                else:
                    if start_digit < end:
                        step = int(charnum(d[start_digit:end]))
                        if (step < 1) or (step > n):
                            jumpout = "yes"
                        else:
                            jumpout = "nearly"
                    else:
                        jumpout = "yes"
                if jumpout != "yes":
                    tour.append(step)
                    city_count = city_count + 1
            if city_count != n:
                flag = "error: not enough cities in tour (only " + str(city_count) + ", should be " + str(n) + ")"
                trace.write(flag + "\n")
            else:
                column_city = 0  # compute actual tour length
                jumpout = "no"
                while jumpout == "no" and column_city < n - 1:
                    row_city = column_city + 1
                    while jumpout == "no" and row_city < n:
                        if tour[column_city] == tour[row_city]:
                            jumpout = "yes"
                        else:
                            row_city = row_city + 1
                    column_city = column_city + 1
                if jumpout == "yes":
                    flag = "error: repetition in tour (city " + str(tour[row_city]) + ")"
                    trace.write(flag + "\n")
                else:   
                    checklength = 0
                    for j in range(0,n-1):
                        checklength = checklength + matrix[tour[j] - 1][tour[j + 1] - 1]
                    checklength = checklength + matrix[tour[n - 1] - 1][tour[0] - 1]
                    if checklength != length:
                        flag = "error: tour lengths don't match (claimed length: " + str(length) + ", actual length: " + str(checklength) + ")"
                        trace.write(flag + "\n")
                    else:
                        flag = "good (length of tour is " + str(length) + ")"
                        trace.write(flag + "\n")
                        return [length]
                        
    return [0]


import os

interested_dir_list = []  # directories I'm interested in (empty if all)
ignore_dir_list = ['IainsExecutions','cityfiles','LastYearsAnswers','Old stuff']  # directories I'm not interested in

interested_tourfile_dir_list = ['TourfileA','TourfileB','TourfileC']  # tour file directories I'm interested in

directlist = []  # get list of directories
for eachitem in os.listdir('.'):
    if os.path.isdir(eachitem):  # check that an item is a directory
        interested = True
        if interested_dir_list != []:  #look to see if I'm interested in this directory
            interested = False
            num_interested_dir = len(interested_dir_list)
            count = 0
            while count < num_interested_dir and interested == False:
                test_dir = interested_dir_list[count]
                if test_dir == eachitem:
                    interested = True
                count = count + 1
        if interested == True:  # true if this is an interesting directory
            to_ignore = False
            if ignore_dir_list != []:  # look to see if I'm disinterested in this directory
                num_disinterested_dir = len(ignore_dir_list)
                count = 0
                while count < num_disinterested_dir and to_ignore == False:
                    test_dir = ignore_dir_list[count]
                    if test_dir == eachitem:
                        to_ignore = True
                    count = count + 1
            if to_ignore == False:  # false if this is an interesting directory that I'm not disinterested in
                directlist.append(eachitem)  # add to the list of directories
directlist.sort()  # sort list of directories

username_name = {}
if os.path.isfile('usernames.txt'):
    users = open('usernames.txt','r')
    x = users.read(1)
    d = ""
    while x != "":
        if ord(x) >= 32 and ord(x) <= 122:
            d = d + x
        x = users.read(1)
    users.close()
    looking_from = 0
    length = len(d)
    if looking_from < length:
        finished = False
    else:
        finished = True
    while finished == False:
        next_comma = d.find(',',looking_from)
        if next_comma != -1:
            username = d[looking_from:next_comma]
            looking_from = next_comma + 1
            next_comma = d.find(',',looking_from)
            if next_comma != -1:
                user = d[looking_from:next_comma]
                looking_from = next_comma + 1
                username_name[username] = user
            else:
                finished = True
        else:
            finished = True

trace = open('trace.txt','w')
trace.write(time.strftime("%a %m/%d/%y", time.localtime()) + "  " + time.strftime("%I:%M:%S %p", time.localtime()) + "\n")

if directlist == []:  # error if I'm interested in no directories
    print('error: no directories present')
    trace.write('error: no directories present\n')
    trace.close()
elif not os.path.isdir('cityfiles'):  # error if there is no directory 'cityfiles'
    print('error: no city-file directory present')
    trace.write('error: no city-file directory present\n')
    trace.close()
elif interested_tourfile_dir_list == []:  # error: there are no tourfile directories of interest
    print('error: no tour-file directories of interest')
    trace.write('error: no tour-file directories of interest\n')
    trace.close()
                           
else:

    os.chdir('cityfiles')  # move to directory 'cityfiles'
    cityfilelist = []
    for eachitem in os.listdir('.'):  # build list of city files
        if os.path.isfile(eachitem):
            cityfilelist.append(eachitem)
    os.chdir('..')         # move up from directory 'cityfiles'
                           
    if cityfilelist == []:  # error if there are no city files
        print('error: no city-files are present')
        trace.write('error: no city-files are present\n')
        trace.close()
    else:
        cityfilelist.sort()  # sort list of city files
        best_tour = {}
        best_person = {}
        
        for currentdir in directlist:  # check each directory I'm interested in

            persons_tour = {}
                           
            os.chdir(currentdir)       # move to directory 'currentdir'
            if currentdir in username_name:
                print_name = username_name[currentdir]
            else:
                print_name = currentdir
            print('doing ' + print_name + ' ...')
            trace.write('\ndirectory: ' + print_name + ' ...\n')
            
            subdirlist = []
            for eachitem in os.listdir('.'):  # get directories in 'currentdir'
                if os.path.isdir(eachitem):   # check that an item is a directory
                    interested = False
                    num_interested_dir = len(interested_tourfile_dir_list)
                    count = 0
                    while count < num_interested_dir and interested == False:
                        test_dir = interested_tourfile_dir_list[count]
                        if test_dir == eachitem:
                            interested = True
                        count = count + 1
                    if interested == True:  # true if I'm interested in this tour file directory
                        subdirlist.append(eachitem)                          
            if subdirlist == []:
                    trace.write('error: no sets of tour-files are present\n')
            else:
                for eachdir in subdirlist:
                           
                    os.chdir(eachdir)  # move to directory 'eachdir'
                    trace.write(eachdir + '...\n')
                    for city_file in cityfilelist:
                        if not os.path.isfile('tour' + city_file):  # look to see if there is a tour file for a city file
                            trace.write(city_file + ': no tour-file present\n')
                        else:
                            filepath = '../../cityfiles/' + city_file  # set filepath for the city file
                            trace.write(city_file + ': ')
                            eachtour = 'tour' + city_file

                            tour_length = tc(filepath,eachtour)[0]

                            if tour_length != 0:
                                if not city_file in persons_tour:
                                    persons_tour[city_file] = tour_length
                                else:
                                    if persons_tour[city_file] > tour_length:
                                        persons_tour[city_file] = tour_length
                    os.chdir('..')  # move back from directory 'eachdir'

                for city_file in cityfilelist:
                    if city_file in persons_tour:
                        if not city_file in best_tour:
                            best_tour[city_file] = persons_tour[city_file]
                            best_person[city_file] = [currentdir]
                        else:
                            if best_tour[city_file] == persons_tour[city_file]:
                                best_person[city_file].append(currentdir)
                            else:
                                if best_tour[city_file] > persons_tour[city_file]:
                                    best_tour[city_file] = persons_tour[city_file]
                                    best_person[city_file] = [currentdir]
            os.chdir('..')  # move back from directory 'currentdir'
    trace.write("\nBest tours:\n\n")
    for city_file in cityfilelist:
        trace.write(city_file + ": ")
        if city_file in best_tour:
            trace.write("best tour has length " + str(best_tour[city_file]) + "\n")
            trace.write(city_file + " winners :")
            length = len(best_person[city_file])
            if length != 0:
                last_person = best_person[city_file][length - 1]
            for person in best_person[city_file]:
                if person in username_name:
                    trace.write(" " + username_name[person])
                else:
                    trace.write(" " + person)
                if person != last_person:
                    trace.write(',')
            trace.write("\n")
        else:
            trace.write("there is no best tour!\n")
    trace.close()
    print('finished')

                                
                                


