/****************************
Maciej Gębala (CC BY-NC 4.0)
server ver. 0.2
2025-04-13
****************************/
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <arpa/inet.h>

#include "./board.h"

int main(int argc, char *argv[]) {
  int server_socket, player_socket[2];
  char nick[2][10];
  socklen_t player_size;
  struct sockaddr_in server_addr, player_addr1, player_addr2;
  char server_message[16], player_message[16];

  int a, b, first, res = 0;
  int player, moveNo, sockNo, move, msg, end_msg = 0;
  bool isWin, isLose, correct;

  if ( argc != 3 ) {
    printf("Wrong number of arguments\n");
    return -1;
  }

  // Create socket
  server_socket = socket(AF_INET, SOCK_STREAM, 0);
  if ( server_socket < 0 ) {
    printf("Error while creating socket\n");
    return -1;
  }
  printf("Socket created successfully\n");

  // Set port and IP
  server_addr.sin_family = AF_INET;
  server_addr.sin_port = htons(atoi(argv[2]));
  server_addr.sin_addr.s_addr = inet_addr(argv[1]);

  // Bind to the set port and IP
  if ( bind(server_socket, (struct sockaddr*)&server_addr, sizeof(server_addr)) < 0 ) {
    printf("Couldn't bind to the port\n");
    return -1;
  }
  printf("Done with binding\n");

  if ( listen(server_socket, 1) < 0 ) {
    printf("Error while listening\n");
    return -1;
  }
  printf("\nListening for incoming connections.....\n");

  // Accept incoming connection 1
  player_size = sizeof(player_addr1);
  player_socket[0] = accept(server_socket, (struct sockaddr*)&player_addr1, &player_size);
  if ( player_socket[0] < 0 ) {
    printf("Can't accept\n");
        return -1;
  }
  printf("Player 1 connected at IP: %s and port: %i\n", inet_ntoa(player_addr1.sin_addr), ntohs(player_addr1.sin_port));

  // Accept incoming connection 2
  player_size = sizeof(player_addr2);
  player_socket[1] = accept(server_socket, (struct sockaddr*)&player_addr2, &player_size);
  if ( player_socket[1] < 0 ) {
    printf("Can't accept\n");
    return -1;
  }
  printf("Player 2 connected at IP: %s and port: %i\n", inet_ntoa(player_addr2.sin_addr), ntohs(player_addr2.sin_port));

  // Player 1 recognition
  memset(server_message, '\0', sizeof(server_message));
  snprintf(server_message, sizeof(server_message), "%s", "700");
  if ( send(player_socket[0], server_message, strlen(server_message), 0) < 0 ) {
    printf("Can't send\n");
    return -1;
  }
  memset(player_message, '\0', sizeof(player_message));
  if ( recv(player_socket[0], player_message, sizeof(player_message), 0) < 0 ) {
    printf("Couldn't receive\n");
    return -1;
  }
  printf("Message from player 1: %s\n", player_message);

  sscanf(player_message, "%d %s", &a, nick[0]);
  if ( a == 1 )
    first = 0;
  else
    first = 1;

  // Player 2 recognition
  memset(server_message, '\0', sizeof(server_message));
  snprintf(server_message, sizeof(server_message), "%s", "700");
  if ( send(player_socket[1], server_message, strlen(server_message), 0) < 0 ) {
    printf("Can't send\n");
    return -1;
  }
  memset(player_message, '\0', sizeof(player_message));
  if ( recv(player_socket[1], player_message, sizeof(player_message), 0) < 0 ) {
    printf("Couldn't receive\n");
    return -1;
  }
  printf("Message from player 2: %s\n", player_message);

  sscanf(player_message, "%d %s", &b, nick[1]);
  if ( a == b ) {
    printf("No player diversity\n");
    return -1;
  }

  memset(server_message, '\0', sizeof(server_message));
  snprintf(server_message, sizeof(server_message), "%s", "600");
  if ( send(player_socket[first], server_message, strlen(server_message), 0) < 0 ) {
     printf("Can't send\n");
     return -1;
  }

  player = 1;
  moveNo = 0;
  isWin = false;
  isLose = false;
  correct = true;
  sockNo = first;

  setBoard();

  while ( (moveNo < 25) && (!isWin) && (!isLose) && correct ) {
    moveNo++;
    printf("Move no: %d\n", moveNo);
    memset(player_message, '\0', sizeof(player_message));
    if ( recv(player_socket[sockNo], player_message, sizeof(player_message), 0) < 0 ) {
      printf("Couldn't receive\n");
      return -1;
    }
    printf("Player %s message: %s\n", nick[sockNo], player_message);
    sscanf(player_message, "%d", &move);
    correct = setMove(move, player);
    if ( correct ) isWin = winCheck(player);
    if ( correct && !isWin ) isLose = loseCheck(player);

    if ( !correct ) { msg = 400; end_msg = 500; res = (player == 1?2:1); }
    else if ( isWin ) { msg = 200+move; end_msg = 100; res = (player == 1?1:2); }
    else if ( isLose ) { msg = 100+move; end_msg = 200; res = (player == 1?2:1); }
    else if ( moveNo == 25 ) { msg = 300+move; end_msg = 300; res = 0; }
    else msg = move;

    memset(server_message, '\0', sizeof(server_message));
    snprintf(server_message, sizeof(server_message), "%d", msg);
    player = 3-player;
    sockNo = 1-sockNo;
    if ( send(player_socket[sockNo], server_message, strlen(server_message), 0) < 0 ) {
      printf("Can't send\n");
      return -1;
    }
    printf("Server message to player %s: %s\n", nick[sockNo], server_message);
    printBoard();
  }

  sockNo = 1-sockNo;
  memset(server_message, '\0', sizeof(server_message));
  snprintf(server_message, sizeof(server_message), "%d", end_msg);
  if ( send(player_socket[sockNo], server_message, strlen(server_message), 0) < 0 ) {
    printf("Can't send\n");
    return -1;
  }
  printf("Server message to player %s: %s\n", nick[sockNo], server_message);

  printf("%s %s %d\n", nick[first], nick[1-first], res);

  // Closing sockets
  close(player_socket[0]);
  close(player_socket[1]);
  close(server_socket);

  printf("End.\n");

  return 0;
}
