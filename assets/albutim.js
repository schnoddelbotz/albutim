"use strict";

if (typeof albumData == 'undefined') {
  var albumData = {};
}
var VIEWMODE_FOLDER = 0; /* const, it is... */
var VIEWMODE_IMAGE = 1;  /* const, it is... */
var currentView = VIEWMODE_FOLDER;
var currentFolder = "/";
var showImageInfo = false;
var runningFullScreen = false;

$(function () {
  if (serveStatically /* defined in <head> */) {
    console.log('I should get albumData from <script src="...">');
    initAlbum();
  } else {
    $.ajax({
      url: "albumdata.json",
      dataType: "json",
      success: function (response) {
        albumData = response;
        initAlbum();
      }
    });
  }
});

$(document).keydown(function (event) {
  var interestedIn = [13, 27, 37, 38, 39, 70, 73, 83];
  var key = event.which;
  if (interestedIn.indexOf(key) != -1) {
    event.preventDefault();
    switch (key) {
      case 39: // right
        if (currentView == VIEWMODE_IMAGE) {
          document.getElementById("nextimg").click()
        } else {
          // in folder mode, show last image
          var f = findPath(currentFolder);
          var i = getFolderImages(f);
          if (i.length > 0) {
            showImage(currentFolder+'/', basename(i[i.length-1]));
          } // FIXME: else, select next folder (enter onEnter)
        }
        break;
      case 37: // left
        if (currentView == VIEWMODE_IMAGE) {
          document.getElementById("previmg").click()
        } else {
          // in folder mode, show first image
          var f = findPath(currentFolder);
          var i = getFolderImages(f);
          if (i.length > 0) {
            showImage(currentFolder+'/', basename(i[0]));
          } // FIXME: else, select previous folder (enter onEnter)
        }
        break;
      case 38: // up
        if (currentView == VIEWMODE_IMAGE) {
          document.getElementById("headerFolderTitle").click()
        } else {
          if (currentFolder == '/') {
            return;
          } else {
            showFolder(dirname(currentFolder));
          }
        }
        break;
      case 27: // esc
        showFolder('/');
        break;
      case 70: // f
        toggleFullScreen();
        break;
      case 83: // s
        console.log('tbd: unhide search/filter form');
        break;
      case 73: // i
        if (showImageInfo) {
          $("#image-info").hide();
          showImageInfo = false;
        } else {
          $("#image-info").show();
          showImageInfo = true;
        }
        break;
        // fixme: folder up!
    }
    // in browse/folder mode:
    // goto/select prev-thumb next-thumb up-folder top-folder 13=enter-sel-folder/img
  }
});

function hashChanged() {
  goHashUrl();
}

function initAlbum() {
  // console.log(albumData);
  document.title = albumData.title;
  $("#headerAlbumTitle").html(albumData.title);
  goHashUrl();
  $("#info-button").click(function () {
    if (showImageInfo) {
      $("#image-info").hide();
      showImageInfo = false;
    } else {
      $("#image-info").show();
      showImageInfo = true;
    }
  });
}

function goHashUrl() {
  if (window.location.hash) {
    var hash = window.location.hash.substr(1);
    // fixme test if exists, otherwise...
    if (hash.match(/\.(jpg|png|gif|mp4)$/i)) {
      var dir = dirname(hash) + '/';
      var img = basename(hash);
      if (dir == '//') {
        dir = '/';
      }
      showImage(dir, img);
    } else {
      showFolder(hash);
    }
  } else {
    window.location.href = '#/';
    showFolder("/");
  }
}

function getSubFolders(folderData) {
  var subFolders = [];
  for (var key in folderData.children) {
    if (folderData.children[key].is_dir) {
      subFolders.push(folderData.children[key].path);
    }
  }
  return subFolders.sort(alphanum);
}

function getFolderImages(folderData) {
  var images = [];
  for (var key in folderData.children) {
    if (folderData.children[key].is_image) {
      images.push(folderData.children[key].path);
    }
  }
  return images.sort(alphanum);
}

function showFolder(imgFolder) {
  currentView = VIEWMODE_FOLDER;
  window.location.href = '#' + imgFolder;
  imgFolder = decodeURIComponent(imgFolder);

  var folder = findPath(imgFolder);
  if (folder == null) {
    show404('folder', imgFolder);
    return;
  }
  currentFolder = imgFolder;

  $("#subHeader").html(albumData.subTitle);
  $("#image-info").hide();
  $("#info-button").hide();
  $("#headerFolderTitle").html(imgFolder);
  $("#headerFolderTitle").attr('href', '#' + imgFolder)

  $("#single-image-container").hide();
  $("#stepnav-thumb-container").hide();
  $("#thumb-container").fadeIn();

  var growingDelay = 0;
  $("#thumb-navigation").html('');

  // always show mouse cursor in folder view
  $('body').css('cursor', 'auto');

  // add folders
  var subFolders = getSubFolders(folder);
  if (subFolders != null) {
    //subFolders = subFolders.sort(alphanum);
    for (var f in subFolders) {
      var folderImage = getFolderImage(subFolders[f]);
      if (folderImage == null) {
        folderImage = 'noFolderPic.jpg';
      }
      var sfImage = 'thumbs' + folderImage;
      var link = '<div class="singlethumb-container float-left">';
      link = link + '<a href="#' + subFolders[f] + '" class="thumb-container">';
      link = link + '<img class="border" src="' + sfImage + '" style="">';
      link = link + '<span class="label">' + subFolders[f] + '</span>';
      link = link + '</a>';
      link = link + '</div>';
      $("#thumb-navigation").append($(link).hide().fadeIn(++growingDelay * 50));
    }
  }

  // add images
  var folderItems = getFolderImages(folder);
  if (folderItems.length > 0) {
    //var keysSorted = folderItems.sort(alphanum);
    for (var i = 0; i < folderItems.length; i++) {
      var imageName = folderItems[i];
      var thumb = '<div class="singlethumb-container">';
      thumb = thumb + '<a href="#' + imageName
          + '" class="thumb-container">';
      if (imageName.match(/\.mp4$/i)) {
        imageName = imageName + '.jpg';
      }
      thumb = thumb + '<img class="border" src="thumbs' + imageName
          + '">';
      thumb = thumb + '</a>';
      thumb = thumb + '</div>';
      $("#thumb-navigation").append(
          $(thumb).hide().fadeIn(++growingDelay * 70));
    }
  }

  return;
  // FIXME:
  // set background image = random image
  var bgImage = 'preview' + getFolderImage(imgFolder);
  $("#background-image").css('background-image', "url('" + bgImage + "')");
}

// get a random image from given folder
function getFolderImage(folder) {
  var folderData = findPath(folder);

  var folderImages = [];
  for (var key in folderData.children) {
    if (folderData.children[key].is_image) {
      folderImages.push(folderData.children[key].path);
    }
  }

  if (folderImages.length > 0) {
    var picNum = Math.floor((Math.random() * folderImages.length) + 0);
    return folderImages[picNum];
  }

  return null;
}

function showImage(folder, image) {
  currentView = VIEWMODE_IMAGE;
  window.location.href = '#' + folder + image;
  folder = decodeURIComponent(folder);
  image = decodeURIComponent(image);
  var folderWithoutTrailingSlash = folder.replace(/\/+$/, "");
  currentFolder = folderWithoutTrailingSlash;
  $("#headerFolderTitle").html(folderWithoutTrailingSlash);
  $("#headerFolderTitle").attr('href', '#' + folderWithoutTrailingSlash);
  $("#single-image-container").show();
  $("#stepnav-thumb-container").show();
  $("#thumb-container").hide();
  $("#subHeader").show();
  if (showImageInfo) {
    $("#image-info").show();
  } else {
    $("#image-info").hide();
  }

  $("#info-button").show();
  $("#subHeader").html(image); // show filename as page sub-title
  var imgdata = findPath(folder + image);
  if (imgdata == null) {
    return show404('image', image);
  }

  if (runningFullScreen) {
    $('body').css('cursor', 'none');
  } else {
    $('body').css('cursor', 'auto');
  }

  // update EXIF info panel
  updateImageInfo(imgdata.exifdata);

  var frontImageUrl = 'preview' + folder + image;
  var backdropImageUrl = 'preview' + folder + image;

  // FIXME: videos/gifs ...
  if (frontImageUrl.match(/.mp4$/i)) {
    frontImageUrl = frontImageUrl + '.jpg';
    backdropImageUrl = backdropImageUrl + '.jpg';
  }
  if (frontImageUrl.match(/.gif$/i)) {
    frontImageUrl = 'original' + folder + image;
  }
  $("#single-image-link").attr('href', 'originals' + folder + image);

  // set background image = current image; fades via css3 transition
  $("#background-image").css('background-image',
      "url('" + backdropImageUrl + "')");
  $("#single-image").attr('src', frontImageUrl);

  // update next/previous buttons
  var folderData = findPath(folderWithoutTrailingSlash);
  var images = getFolderImages(folderData);
  //images = images.sort(alphanum);
  var currentImageId = images.indexOf(folder + image);
  var previous = currentImageId > 0 ? images[currentImageId - 1] : false;
  var next = currentImageId < images.length ? images[currentImageId + 1] : false;
  updatePrevNext(folderWithoutTrailingSlash, previous, next);
}

function updatePrevNext(folder, prev, next) {
  // fixme mp4/gif
  if (next) {
    $("#nextimg").attr('href', '#' + next);
    if (next.match(/\.mp4$/i)) {
      next = next + '.jpg';
    }
    $("#nimg").attr('src', 'thumbs' + next);
  } else {
    $("#nextimg").attr('href', '#' + folder);
    $("#nimg").attr('src', 'assets/folder-up.svg'); // FIXME getImageForFolder
  }
  if (prev) {
    $("#previmg").attr('href', '#' + prev);
    if (prev.match(/\.mp4$/i)) {
      prev = prev + '.jpg';
    }
    $("#pimg").attr('src', 'thumbs' + prev);
  } else {
    $("#previmg").attr('href', '#' + folder);
    $("#pimg").attr('src', 'assets/folder-up.svg'); // FIXME
  }
}

function updateImageInfo(data) {
  var fields = ['dateTime', 'model', 'fNum', 'exposureMode', 'flash',
    'exposureTime', 'exposureBiasValue', 'ISOSpeedRatings', 'focalLength',
    'focalLengthIn35mmFilm', 'whiteBalance', 'exifImageLength',
    'exifImageWidth', 'fileSize'];
  for (var m in fields) {
    $('#' + fields[m]).html("n/a");
  }
  for (var m in data) {
    $('#' + m).html(data[m]);
  }
}

function show404(type, path) {
  // FIXME make it more prominent!
  $("#subHeader").html(
      '<span style="color:red">' + type + ' ' + path + ' not found!'
      + "</span>");
}

/****************** helper functions below ***********************************/

// https://stackoverflow.com/questions/9133500/how-to-find-a-node-in-a-tree-with-javascript
function findPath(needle) {
  var stack = [], node, ii;
  stack.push(albumData.data);
  while (stack.length > 0) {
    node = stack.pop();
    if (node.path == needle) {
      // Found it!
      return node;
    } else if (node.children && node.children.length) {
      for (ii = 0; ii < node.children.length; ii += 1) {
        stack.push(node.children[ii]);
      }
    }
  }
  return null;
}

// http://stackoverflow.com/questions/2532218/pick-random-property-from-a-javascript-object
function pickRandomProperty(obj) {
  var result;
  var count = 0;
  for (var prop in obj) {
    if (Math.random() < 1 / ++count) {
      result = prop;
    }
  }
  return result;
}

// http://stackoverflow.com/questions/13303151/getting-fullscreen-mode-to-my-browser-using-jquery
function toggleFullScreen() {
  if (!document.fullscreenElement &&    // alternative standard method
      !document.mozFullScreenElement && !document.webkitFullscreenElement
      && !document.msFullscreenElement) {  // current working methods
    if (document.documentElement.requestFullscreen) {
      document.documentElement.requestFullscreen();
    } else if (document.documentElement.msRequestFullscreen) {
      document.documentElement.msRequestFullscreen();
    } else if (document.documentElement.mozRequestFullScreen) {
      document.documentElement.mozRequestFullScreen();
    } else if (document.documentElement.webkitRequestFullscreen) {
      document.documentElement.webkitRequestFullscreen(
          Element.ALLOW_KEYBOARD_INPUT);
    }
    runningFullScreen = true;
  } else {
    if (document.exitFullscreen) {
      document.exitFullscreen();
    } else if (document.msExitFullscreen) {
      document.msExitFullscreen();
    } else if (document.mozCancelFullScreen) {
      document.mozCancelFullScreen();
    } else if (document.webkitExitFullscreen) {
      document.webkitExitFullscreen();
    }
    runningFullScreen = false;
    $('body').css('cursor', 'auto');
  }
}

// http://www.davekoelle.com/files/alphanum.js
function alphanum(a, b) {
  function chunkify(t) {
    var tz = new Array();
    var x = 0, y = -1, n = 0, i, j;
    while (i = (j = t.charAt(x++)).charCodeAt(0)) {
      var m = (i == 46 || (i >= 48 && i <= 57));
      if (m !== n) {
        tz[++y] = "";
        n = m;
      }
      tz[y] += j;
    }
    return tz;
  }

  var aa = chunkify(a);
  var bb = chunkify(b);
  for (var x = 0; aa[x] && bb[x]; x++) {
    if (aa[x] !== bb[x]) {
      var c = Number(aa[x]), d = Number(bb[x]);
      if (c == aa[x] && d == bb[x]) {
        return c - d;
      } else {
        return (aa[x] > bb[x]) ? 1 : -1;
      }
    }
  }
  return aa.length - bb.length;
}

// https://github.com/os-js/OS.js/blob/30c7ccac294936dcd2a62c640e4759edb524841e/src/client/javascript/utils/fs.js
function dirname(f) {
  var pstr = f.split(/^(.*)\:\/\/(.*)/).filter(function (n) {
    return n !== '';
  });
  var args = pstr.pop();
  var prot = pstr.pop();
  var result = '';
  var tmp = args.split('/').filter(function (n) {
    return n !== '';
  });
  if (tmp.length) {
    tmp.pop();
  }
  result = tmp.join('/');
  if (!result.match(/^\//)) {
    result = '/' + result;
  }
  if (prot) {
    result = prot + '://' + result;
  }
  return result;
}

function basename(str) {
  return new String(str).substring(str.lastIndexOf('/') + 1);
}
