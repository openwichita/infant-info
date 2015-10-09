(function (window, document) {
  var deleteIcons = document.getElementsByClassName("delete-admin-user"),
      editIcons = document.getElementsByClassName("edit-admin-user");
  for(var i = 0; i < deleteIcons.length; i++) {
    deleteIcons[i].onclick = function(e) {
      var userName = this.parentElement.parentElement.getAttribute("data-user");
      location.href = "/admin/users/delete/"+userName;
    };
  }
  for(var i = 0; i < editIcons.length; i++) {
    editIcons[i].onclick = function(e) {
      var userName = this.parentElement.parentElement.getAttribute("data-user");
      location.href = "/admin/users/edit/"+userName;
    };
  }
}(this, this.document));
