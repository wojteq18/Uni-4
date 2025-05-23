/****************************
Maciej Gębala (CC BY-NC 4.0)
Random Boot ver. 0.2
2025-04-13
****************************/
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <gsl/gsl_rng.h>
#include <gsl/gsl_randist.h>
#include <time.h>
#include <arpa/inet.h>

#include "./board.h"

int randomMove(int player, gsl_rng *generator) {
  int n = 0;
  for (int i=0; i < 5; i++)
    for (int j=0; j < 5; j++)
      if (board[i][j] == 0)
        n++;
  n = gsl_rng_uniform_int(generator, n	);
  for (int i=0; i < 5; i++)
    for (int j=0; j < 5; j++)
      if (board[i][j] == 0)
        if (n == 0)
          return (i+1)*10+j+1;
        else
          n--;
}

int main(int argc, char *argv[]) {
  int server_socket;
  struct sockaddr_in server_addr;
  char server_message[16], player_message[16];

  bool end_game;
  int player, msg, move;

  gsl_rng *generator = gsl_rng_alloc(gsl_rng_mt19937);
  gsl_rng_set(generator, time(NULL));

  if ( argc != 5 ) {
    printf("Wrong number of arguments\n");
    return -1;
  }

  // Create socket
  server_socket = socket(AF_INET, SOCK_STREAM, 0);
  if ( server_socket < 0 ) {
    printf("Unable to create socket\n");
    return -1;
  }
  printf("Socket created successfully\n");

  // Set port and IP the same as server-side
  server_addr.sin_family = AF_INET;
  server_addr.sin_port = htons(atoi(argv[2]));
  server_addr.sin_addr.s_addr = inet_addr(argv[1]);

  // Send connection request to server
  if ( connect(server_socket, (struct sockaddr*)&server_addr, sizeof(server_addr)) < 0 ) {
    printf("Unable to connect\n");
    return -1;
  }
  printf("Connected with server successfully\n");

  // Receive the server message
  memset(server_message, '\0', sizeof(server_message));
  if ( recv(server_socket, server_message, sizeof(server_message), 0) < 0 ) {
    printf("Error while receiving server's message\n");
    return -1;
  }

  memset(player_message, '\0', sizeof(player_message));
  snprintf(player_message, sizeof(player_message), "%s %s", argv[3], argv[4]);
  // Send the message to server
  if ( send(server_socket, player_message, strlen(player_message), 0) < 0 ) {
    printf("Unable to send message\n");
    return -1;
  }

  setBoard();
  end_game = false;
  sscanf(argv[3], "%d", &player);

  while ( !end_game ) {
    memset(server_message, '\0', sizeof(server_message));
    if ( recv(server_socket, server_message, sizeof(server_message), 0) < 0 ) {
      printf("Error while receiving server's message\n");
      return -1;
    }
    sscanf(server_message, "%d", &msg);
    move = msg%100;
    msg = msg/100;
    if ( move != 0 ) {
      setMove(move, 3-player);
    }
    if ( (msg == 0) || (msg == 6) ) {
      move = randomMove(player, generator);
      setMove(move, player);
      memset(player_message, '\0', sizeof(player_message));
      snprintf(player_message, sizeof(player_message), "%d", move);
      if ( send(server_socket, player_message, strlen(player_message), 0) < 0 ) {
        printf("Unable to send message\n");
        return -1;
      }
     } else {
       end_game = true;
       switch ( msg ) {
         case 1 : printf("You won.\n"); break;
         case 2 : printf("You lost.\n"); break;
         case 3 : printf("Draw.\n"); break;
         case 4 : printf("You won. Opponent error.\n"); break;
         case 5 : printf("You lost. Your error.\n"); break;
      }
    }
  }

  // Close socket
  close(server_socket);

  return 0;
}
