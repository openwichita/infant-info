(function (window, document) {
  var deleteIcons = document.getElementsByClassName("delete-admin-user"),
      editIcons = document.getElementsByClassName("edit-admin-user"),
      addNewButton = document.getElementById("addUserButton");
  addNewButton.onclick = function(e) {
    location.href = "/admin/users/create";
  }
  for(var i = 0; i < deleteIcons.length; i++) {
    deleteIcons[i].onclick = function(e) {
      var userName = this.parentElement.parentElement.getAttribute("data-user");
      var answer = confirm("Are you sure you want to delete user '"+userName+"'?");
      if(answer) {
        location.href = "/admin/users/delete/"+userName;
      }
    };
  }
  for(var i = 0; i < editIcons.length; i++) {
    editIcons[i].onclick = function(e) {
      var userName = this.parentElement.parentElement.getAttribute("data-user");
      location.href = "/admin/users/edit/"+userName;
    };
  }
}(this, this.document));
