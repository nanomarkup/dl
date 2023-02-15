# Dependency Langauge
The Dependency Language provides methods for managing dependencies between items and store data into the file.

The first line should be the same as file extension. For example, it should be "sb" for "apps.sb" file. <sub>It should be removed.</sub>
#### Item
```
path.item:
```
- Item name is a full path to the item
- Item name starts from a new line
- It does not include spaces 
- It ends with the ":" token
#### Field
```
path.item:
    field
```
- Field name starts from a new line
- It does not include spaces 
- It can be defined after item declaration
#### Value
```
path.item:
    field 
    field 1
    field "text"
    field path.item2
    field *path.item2
    field path.function(1, "text")
```
- It can be empty
- It can be a number
- It can be a text enclosed in quotation marks
- It can be an item name
- It can be a reference to the item using "*" as a prefix for item name
- It can be a function with/out parameters
#### Group item
```
path.item:
    field [group]path.writer
  
path.writer:
    Message "Hello"
	
[group]path.writer:
    Message "Hi"
```
- Group name enclosed in square brackets
- It locates before item name
- It gives the ability to overload the initialization of the same item using different values
#### Defines
```
defines:
    name github.com/sapplications
  
{name}/dl:
    field {name}/sb
```
- It is a special "defines" item
- The field name is a define name
- The value cannot be empty
- The define name enclosed in curly brackets will be replaced by the define value 
#### Internal Initialization
```
path.item:
    field path.writer {
        Message "Hi"
    }
  
path.writer:
    Message "Hello"
```
- Use curly brackets for internal item initialization
- It gives the ability to overload the initialization of the same item using different values
- The left curly bracket locats on the same line as the item name (at the end)
- The right curly bracket starts from a new line
#### Comments
```
// This is a comment
```
- It starts from a new line
- It begins with double slashes
