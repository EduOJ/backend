

# License.

This project is licensed under [GNU AFFERO GENERAL PUBLIC LICENSE Version 3](./license.md).

# Contribution

Our document is still under construction.

# Code style.
All code must be formatted by go fmt.

All tests of app/controllers is running with the same in-memory database,
so they shouldn't rely on a clean database and shouldn't cleanup after runs.
Also, they should be running under parallel mode.

All other tests should make it own database and clean it up after running.

# Permissions

Here are the permissions and their descriptions.

|      Name      |                     Description                     |
|:--------------:|:---------------------------------------------------:|
|   creat_user   |           the permission to create users            |
|   update_user  |           the permission to update users            |
|   delete_user  |           the permission to delete users            |
|    get_user    | the permission to get all the information of a user |
|   get_users    |            the permission to get users              |