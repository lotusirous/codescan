# Code scan

A simple static code analysis

[API document](./doc/api.md)

[Development document](./doc/development.md)

[Scanner design](./doc/scanner-design.md)

[Design document](./doc/scanner-design.md)

# Challenge

Requirement:
Change your codes to do finding deduplication
i.e. do not add new findings to DB for existing findings that report as the same with previous scans.
Example:

* At the first scan, we only have file A with the line content:
        ...
  5     private_key = "123456"
        ...

 ==> One finding found: file_A@5

* At the second scan:
  File A does not change
  File B is added new:

     ...
  10    private_key = "abcd"
     ...
==> Two findings found:
  (1)  file_A@5
  (2)  file_B@10

 For this case, we only need update timestamp for finding 1, but add new for finding 2.
