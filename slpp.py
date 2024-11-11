# originally by chp2001 on the Reassembly Discord


import re
import math
CONST_quotes = '"\'['
CONST_num = set('*/+-eExX.ABCDEFabcdef0123456789')

class ParseError(Exception):
    pass


class SLPP:

    def __init__(self):
        self.text = ''
        self.ch = ''
        self.at = 0
        self.len = 0
        self.depth = 0
        self.newline = '\n'
        self.tab = '  '
        self.openbrackets = 0
        self.closebrackets = 0

    def decode(self, text):
        if not text or not isinstance(text, str):
            return
        self.text = text
        self.at, self.ch, self.depth = 0, '', 0
        self.len = len(text)
        self.next_chr()
        result = self.value()
        return result

    def encode(self, obj):
        if not obj:
            return
        self.depth = 0
        return self.__encode(obj)

    def __encode(self, obj):
        s = ''
        tab = self.tab
        newline = self.newline
        tp = type(obj)
        if isinstance(obj, str):
            s += '"%s"' % obj.replace(r'"', r'\"')
        elif tp in [int, float, int, complex]:
            s += str(obj)
        elif tp is bool:
            s += str(obj).lower()
        elif tp in [list, tuple, dict]:
            self.depth += 1
            if len(obj) == 0 or ( tp is not dict and len(list(filter(
                    lambda x:  type(x) in (int,  float,  int) \
                    or (isinstance(x, str) and len(x) < 10),  obj
                ))) == len(obj) ):
                newline = tab = ''
            dp = tab * self.depth
            sep = '%s%s' % (newline, dp) if len(obj) > 4 else ""
            s += sep + "{"
            js = ', ' + sep
            if tp is dict:
                s += js.join([self.__encode(v) if type(k) is int \
                              else '%s = %s' % (k, self.__encode(v)) \
                              for k, v in obj.items()
                          ])
            else:
                s += js.join([self.__encode(el) for el in obj])
            self.depth -= 1
            s += "}"
        elif tp is set:
            return "|".join(obj)
        return s

    def white(self):
        in_comment = False
        while self.ch:
            if (self.ch == '-' and self.text[self.at] == '-') or self.ch == '#':
                in_comment = True
            if in_comment or self.ch.isspace():
                if self.ch == '\n':
                    in_comment = False
                self.next_chr()
            else:
                break

    def next_chr(self):
        if self.at >= self.len:
            self.ch = None
            return None
        self.ch = self.text[self.at]
        self.at += 1
        return True

    def value(self):
        self.white()
        if not self.ch:
            return
        elif self.ch == '{':
            self.openbrackets+=1
            return self.object()
        elif (self.ch=='-' and self.text[self.at]=='-'):
            self.white()
        if self.ch == "[":
            self.next_chr()
        if self.ch.isdigit() or self.ch == '-' or (self.ch=='p' and self.text[self.at]=='i') or (self.ch=='i' and self.text[self.at]=='n' and self.text[self.at+1]=='f'):
            return self.number()
        elif self.ch in CONST_quotes:
            return self.string()
        return self.word()

    def string(self):
        s = u''
        start = self.ch
        end = ']' if start == '[' else start
        while self.next_chr():
            if self.ch == end:
                self.next_chr()
                return s
            elif self.ch == '\\':
                self.next_chr()
                if self.ch == "\n":
                    continue
                elif self.ch == "n":
                    s += '\n'
                    continue
                elif self.ch != end:
                    s += '\\'
            s += self.ch
        raise ParseError(u'Unexpected end of string while parsing Lua string: %s%s%s' \
                         % (start, s[:20].replace("\n", "\\n") + ("..." if len(s) >= 20 else ""), end))

    def object(self):
        o = {}
        idx = 0
        numeric_keys = False
        self.next_chr()
        while self.ch and self.ch != '}':
            self.white()
            if self.ch=='}':
                continue
            k = self.value()
            if self.ch == ']':
                numeric_keys = True
                self.next_chr()
            self.white()
            if self.ch == '=':
                self.next_chr()
                self.white()
                o[k] = self.value()
            else:
                o[idx] = k
            idx += 1
            self.white()
            if self.ch == ',': # optional comma
                self.next_chr()
            self.white()
        if self.ch != '}':
            leadup=''
            for num in range(-40,0):
                leadup+=self.text[self.at+num]
            # b1=0
            # b2=0
            # for char in range(0,self.at):
            #     if self.text[char]=='{':
            #         b1+=1
            #     if self.text[char]=='}':
            #         b2+=1
            raise ParseError('Unexpected end of data while parsing Lua table at '+str(self.at)+' leadup was '+leadup+' bracket comparison is '+str(self.openbrackets)+' '+str(self.closebrackets)+' total number of characters '+str(len(self.text)))
        else:
            self.closebrackets+=1
        # b1=0
        # b2=0
        # for char in range(0,self.at):
        #     if self.text[char]=='{':
        #         b1+=1
        #     if self.text[char]=='}':
        #         b2+=1
        # if b1==self.openbrackets and b2!=self.closebrackets:
        #     leadup=''
        #     for num in range(-40,0):
        #         leadup+=self.text[self.at+num]
        #     raise ParseError('Uncaught bracket at '+str(self.at)+' leadup was '+leadup+' bracket comparison is '+str(self.openbrackets)+' '+str(self.closebrackets)+' total number of characters '+str(len(self.text))+' bracket comparison is '+str(b1)+' '+str(b2))
        self.next_chr()
        if not numeric_keys and all(isinstance(key, int) for key in o):
            ar = []
            for key in o:
               ar.insert(key, o[key])
            o = ar
        return o

    def word(self):
        s = ''
        if self.ch != '\n':
            s = self.ch
        while self.next_chr():
            if (self.ch.isalnum() or self.ch in "_-|") and not (self.ch=='-' and self.text[self.at]=='-') :
                s += self.ch
            else:
                if re.match('^true$', s, re.I):
                    return True
                elif re.match('^false$', s, re.I):
                    return False
                elif s == 'nil':
                    return None
                elif '|' in s:  # special bitset syntax AA|BB|CC
                    return set(s.split('|'))
                return str(s)

    def operation(self,vals,index,operator):
        newvals=[]
        val1=0
        val2=0
        #print('vals contains '+str(vals))
        #print('index is '+str(index))
        #print('operating on '+vals[index-1]+' '+vals[index]+' '+vals[index+1])
        if index>2:
            for num in range(0,index-1):
                newvals.append(vals[num])
        if '.' in vals[index-1]:
            val1=float(vals[index-1])
        else:
            val1=int(vals[index-1])
        if '.' in vals[index+1]:
            val2=float(vals[index+1])
        else:
            val2=int(vals[index+1])
        if operator=='e':
            newvals.append(str(val1*10**val2))
        elif operator=='*':
            newvals.append(str(val1*val2))
        elif operator=='/':
            newvals.append(str(val1/val2))
        elif operator=='+':
            newvals.append(str(val1+val2))
        elif operator=='-':
            newvals.append(str(val1-val2))
        if len(vals)-1>index+2:
            for num in range(index+2,len(vals)):
                newvals.append(vals[num])
        return newvals

    def equation(self,n):
        changing=True
        vals=[]
        vals.append(n)
        while changing:
            changing=False
            val=[]
            indexnum=0
            for key in vals:
                if 'e' in key and key!='e':
                    splitter='e'
                    splits=key.split(splitter)
                    changing=True
                    val.append(splits[0])
                    val.append(splitter)
                    if len(splits)>2:
                        for num in range(1,len(splits)-1):
                            val.append(splits[num])
                            val.append(splitter)
                    val.append(splits[len(splits)-1])
                    break
                if '*' in key and key!='*':
                    splitter='*'
                    splits=key.split(splitter)
                    changing=True
                    val.append(splits[0])
                    val.append(splitter)
                    if len(splits)>2:
                        for num in range(1,len(splits)-1):
                            val.append(splits[num])
                            val.append(splitter)
                    val.append(splits[len(splits)-1])
                    break
                if '/' in key and key!='/':
                    splitter='/'
                    splits=key.split(splitter)
                    changing=True
                    val.append(splits[0])
                    val.append(splitter)
                    if len(splits)>2:
                        for num in range(1,len(splits)-1):
                            val.append(splits[num])
                            val.append(splitter)
                    val.append(splits[len(splits)-1])
                    break
                if '+' in key and key!='+':
                    splitter='+'
                    splits=key.split(splitter)
                    changing=True
                    val.append(splits[0])
                    val.append(splitter)
                    if len(splits)>2:
                        for num in range(1,len(splits)-1):
                            val.append(splits[num])
                            val.append(splitter)
                    val.append(splits[len(splits)-1])
                    break
                if '-' in key and key!='-' and indexnum>1 and not vals[indexnum-1] in 'e*/+-':
                    splitter='-'
                    splits=key.split(splitter)
                    changing=True
                    val.append(splits[0])
                    val.append(splitter)
                    if len(splits)>2:
                        for num in range(1,len(splits)-1):
                            val.append(splits[num])
                            val.append(splitter)
                    val.append(splits[len(splits)-1])
                    break
                indexnum+=1
            if changing==True:
                #print('removing '+vals[indexnum])
                vals.remove(vals[indexnum])
                #print('vals contains '+str(vals))
                for num in range(len(val)-1,-1,-1):
                    if val[num]=='':
                        continue
                    vals.insert(indexnum,val[num])
        indexnum=0
        changing=True
        while changing:
            case=None
            changing=False
            indexnum=0
            for key in vals:
                if key=='e':
                    case=key
                    break
                if key=='*':
                    case=key
                    break
                if key=='/':
                    case=key
                    break
                if key=='+':
                    case=key
                    break
                if key=='-' and indexnum>1 and not vals[indexnum-1] in 'e*/+-':
                    case=key
                    break
                indexnum+=1
            if case!=None:
                vals=self.operation(vals,indexnum,case)
                changing=True
        if len(vals)==1:
            return vals[0]

    def globalvar(self):
        val=0
        if self.ch=='i':
            val=2**31
        elif self.ch=='p':
            val=math.pi
        while self.ch:
            if self.ch in 'infp':
                self.next_chr()
            else:
                break
        return str(val)

    def number(self):
        n = ''
        operate=False
        while self.ch:
            if self.ch in 'e*/+-':
                if self.ch=='-' and self.text[self.at]=='-':
                    break
                #print('found equation!')
                n += self.ch
                self.next_chr()
                #operate=True
            elif self.ch in CONST_num:
                n += self.ch
                self.next_chr()
            elif self.ch in 'infp':
                n += self.globalvar()
            else:
                break
        color=False
        if 'x' in n:
            color=True
        elif '*' in n or '/' in n or '+' in n or '-' in n or 'e' in n:
            operate=True
        if operate:
            n=self.equation(n)
        if not color and ("." in n or 'e' in n or float(n)<1) :
            #print(n)
            if float(n)-int(float(n))==0:
                return int(float((n))
            return float(n)
        else:
            #print('returning '+n)
            return int(n, 0)

slpp = SLPP()

__all__ = ['slpp']
