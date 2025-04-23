 ____ ____ ____ ____ ____ ____ ____ _________ ____ ____ ____ ____ ____ ____ ____ 
||S |||e |||a |||r |||c |||h |||X |||       |||I |||n |||d |||e |||x |||e |||r ||
||__|||__|||__|||__|||__|||__|||__|||_______|||__|||__|||__|||__|||__|||__|||__||  
|/__\|/__\|/__\|/__\|/__\|/__\|/__\|/_______\|/__\|/__\|/__\|/__\|/__\|/__\|/__\|  

===================================================================================
                                  SearchX Indexer
===================================================================================

This repository is part of the **SearchX** project. The main objective of this
module is to process **data leaks** (data leaks) and store them securely and
efficiently in a **SQL** database.

With **SearchX Indexer**, you can ensure that sensitive leak data is processed,
organized and saved in an optimized way, allowing easy integration with other
modules of the **SearchX** project.

-----------------------------------------------------------------------------------
                                    Features
-----------------------------------------------------------------------------------

- **Efficient Leak Processing**: The module is designed to process large volumes
of data quickly, ensuring that you can deal with data leaks efficiently.

- **Secure Storage**: Uses best practices to store sensitive information
in a SQL database in a secure manner, encrypting or masking data when necessary.

- **Simple Integration**: Easy integration with other modules of the**SearchX**
project, allowing the construction of more robust and scalable solutions.

- **Scalability**: The module's design allows it to support large volumes of data,
being able to deal with large-scale leaks.

- **Monitoring and Logs**: Built-in monitoring tools to ensure indexing and
storage processes run smoothly, with detailed logs for auditing and debugging.

-----------------------------------------------------------------------------------
                                  Prerequisites
-----------------------------------------------------------------------------------

- Go 1.23 or higher.
- PostgreSQL database.
- Network access to integrate with other **SearchX** modules.

-----------------------------------------------------------------------------------
                                     Build
-----------------------------------------------------------------------------------

mkdir -p $HOME/go/src/github.com/0x53bin/ \
&& git clone https://github.com/0x53binx/searchx-indexer \ 
$HOME/go/src/github.com/0x53bin/searchx-indexer \
&& cd $HOME/go/src/github.com/0x53bin/searchx-indexer \
&& go build cmd/main.go -o searchx-indexer
