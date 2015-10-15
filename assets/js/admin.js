(function (window, document) {
  var deleteUserIcons = document.getElementsByClassName("delete-admin-user"),
      editUserIcons = document.getElementsByClassName("edit-admin-user"),
      addNewUserButton = document.getElementById("addUserButton"),
      deleteResourceIcons = document.getElementsByClassName("delete-resource"),
      editResourceIcons = document.getElementsByClassName("edit-resource"),
      addNewResourceButton = document.getElementById("addResourceButton");
  /* User Management */
  if(addNewUserButton) {
    addNewUserButton.onclick = function(e) {
      location.href = "/admin/users/create";
    }
    for(var i = 0; i < deleteUserIcons.length; i++) {
      deleteUserIcons[i].onclick = function(e) {
        var userName = this.parentElement.parentElement.getAttribute("data-user");
        var answer = confirm("Are you sure you want to delete user '"+userName+"'?");
        if(answer) {
          location.href = "/admin/users/delete/"+encodeURIComponent(userName);
        }
      };
    }
    for(var i = 0; i < editUserIcons.length; i++) {
      editUserIcons[i].onclick = function(e) {
        var userName = this.parentElement.parentElement.getAttribute("data-user");
        location.href = "/admin/users/edit/"+encodeURIComponent(userName);
      };
    }
  }

  /* Resource Management */
  if(addNewResourceButton) {
    addNewResourceButton.onclick = function(e) {
      location.href = "/admin/resources/create";
    }
    for(var i = 0; i < deleteResourceIcons.length; i++) {
      deleteResourceIcons[i].onclick = function(e) {
        var resTitle = this.parentElement.parentElement.getAttribute("data-resource");
        var answer = confirm("Are you sure you want to delete resource '"+resTitle+"'?");
        if(answer) {
          location.href = "/admin/resources/delete/"+encodeURIComponent(resTitle);
        }
      };
    }
    for(var i = 0; i < editResourceIcons.length; i++) {
      editResourceIcons[i].onclick = function(e) {
        var resTitle = this.parentElement.parentElement.getAttribute("data-resource");
        location.href = "/admin/resources/edit/"+encodeURIComponent(resTitle);
      };
    }
  }
}(this, this.document));
