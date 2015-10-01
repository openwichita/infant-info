(function (window, document) {
  var layout = document.getElementById('layout'),
      menu = document.getElementById('menu'),
      menuLink = document.getElementById('menuLink'),
      flashMsg = document.getElementById('flashmessage');

  function toggleClass(element, className) {
    var classes = element.className.split(/\s+/),
        length = classes.length,
        i = 0;
    for(; i < length; i++) {
      if(classes[i] === className) {
        classes.splice(i, 1);
        break;
      }
    }
    if(length === classes.length) {
      // The className wasn't found
      classes.push(className);
    }
    element.className = classes.join(' ');
  }

  if(flashMsg != null && flashMsg.innerText != "") {
     flashMsg.removeAttribute('hidden');
  }

  menuLink.onclick = function(e) {
    var active = 'active';
    e.preventDefault();
    toggleClass(layout, active);  
    toggleClass(menu, active);
    toggleClass(menuLink, active);
  };
}(this, this.document));
