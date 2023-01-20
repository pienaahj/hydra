/*
$(function() {
    var socket = null;
    var msgBox = $("#chatbox textarea");
    var label = $("#chatbox label");
    var messages = $("#messages");
    $("#chatbox").submit(function() {
      if (!msgBox.val()) return false;
      if (!socket) {
        alert("Error: There is no socket connection.");
        return false;
      }
      socket.send(label.text() + " " + msgBox.val() + "\n");
      msgBox.val("");
      return false;
    });
    if (!window["WebSocket"]) {
      alert("Error: Your browser does not support web sockets.")
    } else {
      if (window.location.protocol == "https:"){
        socket = new WebSocket("wss://localhost:8062/ChatRoom/");
      } else {
        socket = new WebSocket("ws://localhost:8061/ChatRoom/");
      }
      socket.onclose = function() {
        alert("Connection has been closed.");
      }
      socket.onmessage = function(e) {
        messages.append(
          $("<li>").append(
            e.data
          ));
      }
    }
  });
  */

  (function() {
    var socket = null;
    var msgBox = document.querySelector("#chatbox textarea");
    var label = document.querySelector("#chatbox label");
    var messages = document.querySelector("#messages");
    document.querySelector("#chatbox").addEventListener("submit", function(event) {
      event.preventDefault();
      if (!msgBox.value) return false; //there is no message
      if (!socket) { //there is no active socket
        alert("Error: There is no socket connection.");
        return false;
      }
      // make the message and send it
      socket.send(label.textContent + " " + msgBox.value + "\n");
      msgBox.value = ""; //clear the message field
      return false;
    });

    // check if browser supports websockets
    if (!window["WebSocket"]) {
      alert("Error: Your browser does not support web sockets.")
    } else {
    // check the protocol and construct a new websocket
      if (window.location.protocol == "https:"){ 
        socket = new WebSocket("wss://localhost:8062/ChatRoom/");
      } else {
        socket = new WebSocket("ws://localhost:8061/ChatRoom/");
      }
      // notify that the websocket has closed
      socket.onclose = function() {
        alert("Connection has been closed.");
      }
      // make the message and send it
      socket.onmessage = function(e) {
        var node = document.createElement("LI");
        var textnode = document.createTextNode(e.data);
        node.appendChild(textnode);
        messages.appendChild(node);
      }
    }
  })();
  