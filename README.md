# p2p
Light data transfer tool that allows clients behind NAT to communicate directly without(constant) 3rd party server - WIP


POC  - UDP hole punching

    client:
        - register to server
        - wait for peer access point
        - establish connection

    server:
        - wait for client regestration
        - once two clients registered, forward them regestration data


POC results:
    Above business logic is working however, in cases where the client is behind *symetric* NAT,
    the router is blocking the communication.

Options:
    1. Might be possible to guess the next avilable port if there is a common policy for NAT routers.
    Hoever, it will probably will prove to be difficuly in multi-client p2p situation.
    2. Relay the communication through a server, however I will try to avoid this path as much as possible as there should be working solution.
    3. One of the artivle below may provide a solution.





Related articales:

    https://pdos.csail.mit.edu/papers/p2pnat.pdf, might remove requirment of 3rd party server.
    http://www1.cs.columbia.edu/~salman/publications/skype1_4.pdf
    https://arxiv.org/ftp/cs/papers/0412/0412017.pdf
    http://citeseerx.ist.psu.edu/viewdoc/download?doi=10.1.1.455.3700&rep=rep1&type=pdf
    https://www.researchgate.net/publication/224062560_Research_on_Symmetric_NAT_Traversal_in_P2P_applications


    * https://patents.google.com/patent/US9497160B1/en some patent, see that I am not sued by mistake  ^^
