(function (window, document) {
  var deleteIcons = document.getElementsByClassName("delete-admin-user");
  for(var i = 0; i < deleteIcons.length; i++) {
    deleteIcons[i].onclick = function(e) {
      var userName = this.getAttribute("data-user");
      window.location = "/admin/users/delete/"+userName;
    };
  }
}(this, this.document));
