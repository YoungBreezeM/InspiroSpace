
# Easy-Gin

## Overview
Easy-Gin program is a boot strap about gin. Easytting to create a web application. 

## Build

1. Build a binary file

    ```shell
    make
    ```
2. Build done,you will has a server binary

## Generate

1. Generate golang protobuf file

    ```shell
    make gen
    ```


## RUN

1. Normal to start

    ```shell
     easy-gin -c configs
    ```
    You can choose env is dev or prod

    Prod
    ```shell
     easy-gin -c configs
    ```
    Dev
    ```shell
     ENV_CONF="dev" easy-gin -c configs
    ```

