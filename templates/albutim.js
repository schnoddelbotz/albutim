
"use strict";

// comes in via json
var albumImages = [];
var albumTitle = "";
var albumCreated = "";
var albumSubtitle = "";

// get computed / saves state
var showImageInfo = false;
var currentFolder;

$(function() {
  // previously via albumdata.json/xhr, now via <script src="albumdata.js"...>
  albumImages = albumData.images;
  albumTitle = albumData.title;
  albumCreated = albumData.created;
  albumSubtitle = albumData.subtitle;
  initAlbum();
});

$(document).keydown(function( event ) {
  //13=enter // event.which == 38 /*up*/
  var interestedIn = [13,27,37,38,39,70,73,83];
  var key = event.which;
  if (interestedIn.indexOf(key) != -1) {
    event.preventDefault();
    //console.log(event.which);
    // currentFolder
    switch(key) {
      case 39: // right
        document.getElementById("nextimg").click()
        break;
      case 37: // left
        document.getElementById("previmg").click()
        break;
      case 38: // up
        document.getElementById("headerFolderTitle").click()
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
  document.title = albumTitle;
  $("#headerAlbumTitle").html(albumTitle);
  // analyze hashtag to route to image directly?
  goHashUrl();
  $( "#info-button" ).click(function() {
    if (showImageInfo) {
      $("#image-info").hide();
      showImageInfo = false;
    } else {
      $("#image-info").show();
      showImageInfo = true;
      //showImage(folder,next);
    }
  });
}

function goHashUrl() {
  if (window.location.hash) {
    var hash = window.location.hash.substr(1);
    // fixme test if exists, otherwise...
    if (hash.match(/\.(jpg|png|gif|mp4)/i) ) {
      var dir = dirname(hash)+'/';
      var img = basename(hash);
      if (dir=='//') {
        dir='/';
      }
      //alert(dir+" und i:"+img);
      showImage(dir,img);
    } else {
      //console.log("Starting at "+hash);
      showFolder(hash);
    }
  } else {
    //console.log("Starting at root");
    showFolder("/");
  }
}

function getSubFolders(imgFolder) {
  var subFolders = [];
  for (var i in albumImages) {
    //console.log("GOT SUB:" + i);
    //console.log("Has par:" + dirname(i));
    if (  (dirname(i)+'/' == imgFolder || dirname(i) == imgFolder)   && i != "/") {
      subFolders.push(i)
    }
  }
  return subFolders;
}

function showFolder(imgFolder) {
  //console.log("Now showing folder: "+imgFolder);
  window.location.href = '#' + imgFolder;
  imgFolder = decodeURIComponent(imgFolder);

  $("#subHeader").html(albumSubtitle);
  $("#image-info").hide();
  $("#info-button").hide();
  $("#headerFolderTitle").html(imgFolder);
  $("#headerFolderTitle").attr('href','#'+imgFolder)

  $("#single-image-container").hide();
  $("#stepnav-thumb-container").hide();
  $("#thumb-container").fadeIn();

  var newContent = "";
  var growingDelay = 0;
  $("#thumb-navigation").html('');

  // add folders
  var subFolders = getSubFolders(imgFolder).sort(alphanum);
  for (var f in subFolders) {
    //console.log("Link: "+subFolders[f]);
    var sfImage = 'thumb' + getFolderImage(subFolders[f]);
    var link = '<div class="singlethumb-container float-left">';
       link = link+ '<a href="#'+subFolders[f]+'" class="thumb-container">';
       link = link+ '<img class="border" src="'+sfImage+'" style="">';
       link = link+ '<span class="label">'+subFolders[f]+'</span>';
       link = link+ '</a>';
       link = link+ '</div>';
    //newContent = newContent + link;
    $("#thumb-navigation").append($(link).hide().fadeIn(++growingDelay*50));
  }
  //$("#thumb-navigation").html(newContent);
  //$("#thumb-navigation").html(newContent);

  // add images
  //
  var folderItems = albumImages[imgFolder];
  var keysSorted = Object.keys(folderItems).sort(alphanum);
  for (var i = 0; i < keysSorted.length; i++) {
    var imageName = keysSorted[i];
    var thumb = '<div class="singlethumb-container">';
    thumb = thumb+'<a href="#'+imgFolder + imageName+'" class="thumb-container">';
    if (imageName.match(/\.mp4$/i)) {
      imageName = imageName + '.jpg';
    }
    thumb = thumb+'<img class="border" src="thumb'+ imgFolder + imageName  +'">';
    thumb = thumb+'</a>';
    thumb = thumb+'</div>';
    $("#thumb-navigation").append($(thumb).hide().fadeIn(++growingDelay*70));
  }

  // update main content
  //$("#thumb-navigation").html(newContent);
  // set background image = random image
  var bgImage = 'preview'+getFolderImage(imgFolder);
  $("#background-image").css('background-image', "url('"+bgImage+"')");
}

function getFolderImage(folder) {
  //console.log("GET IMAGES FOR: "+folder);
  var imgs = Object.keys(albumImages[folder]);
  if (imgs.length > 0) {
    var picNum = Math.floor((Math.random() * imgs.length) + 0);
    // FIXME whatif this is a folder ... can it be?
    var folderImage = folder+imgs[picNum];
    if (folderImage.match(/\.mp4$/i)) {
      folderImage = folderImage + '.jpg';
    }
    return folderImage;
  } else {
    if (folder=='/') {
      // try harder / use random pic / default pic if root folder / HACK
      var rootFolders = Object.keys(albumImages);
      var imageFound = false;
      var iterations = 0;
      var randImage = "no-folder-icon.jpg"; // --default-image
      var randRootFolder;
      while (!imageFound && iterations<100) {
        randRootFolder = pickRandomProperty(albumImages);
        if (randRootFolder) {
          var rfolder = albumImages[randRootFolder]
          var rimage  = pickRandomProperty(rfolder);
          // fixme check if non-empty
          if (rimage.match(/jpg/i)) {
            randImage = rimage;
            imageFound = true;
          }
        }
        iterations++;
      }
      return randRootFolder + randImage;
    }
    return "no-folder-icon.jpg";
  }
}

function showImage(folder,image) {
  window.location.href = '#' + folder + image;
  folder = decodeURIComponent(folder);
  $("#headerFolderTitle").html(folder);
  $("#headerFolderTitle").attr('href','#'+folder);
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
  var imgdata = albumImages[folder][image];
  updateImageInfo(imgdata);
  var frontImageUrl = 'preview'+folder+image;
  var backdropImageUrl = 'preview'+folder+image;
  if (frontImageUrl.match(/.mp4$/i)) {
    frontImageUrl = frontImageUrl + '.jpg';
    backdropImageUrl = backdropImageUrl + '.jpg';
  }
  if (frontImageUrl.match(/.gif$/i)) {
    frontImageUrl = 'original'+folder+image;
  }
  $("#single-image-link").attr('href','original'+folder+image);

  // set background image = current image; fades via css3 transition
  $("#background-image").css('background-image', "url('"+backdropImageUrl+"')");
  $("#single-image").attr('src', frontImageUrl);

  // update next/previous buttons
  var previous = oFunctions.keys.previous(albumImages[folder], image);
  var next     = oFunctions.keys.next(albumImages[folder], image);
  updatePrevNext(folder, previous, next);
}

function updatePrevNext(folder, prev, next) {
  // fixme mp4/gif
  if (next) {
    $("#nextimg").attr('href', '#'+folder+next);
    if (next.match(/\.mp4$/i)) {
      next = next + '.jpg';
    }
    $("#nimg").attr('src','thumb'+folder+next);
  } else {
    $("#nextimg").attr('href', '#'+folder);
    $("#nimg").attr('src','folder-up.svg'); // FIXME getImageForFolder
  }
  if (prev) {
    $("#previmg").attr('href', '#'+folder+prev);
    if (prev.match(/\.mp4$/i)) {
      prev = prev + '.jpg';
    }
    $("#pimg").attr('src','thumb'+folder+prev);
  } else {
    $("#previmg").attr('href', '#'+folder);
    $("#pimg").attr('src','folder-up.svg'); // FIXME
  }
}

function updateImageInfo(data) {
  // FIXME reset all children of container first...
  for (var m in data) {
    $('#'+m).html(data[m]);
  }
}

/****************** helper functions below ***********************************/

// http://stackoverflow.com/questions/2532218/pick-random-property-from-a-javascript-object
function pickRandomProperty(obj) {
  var result;
  var count = 0;
  for (var prop in obj)
    if (Math.random() < 1/++count)
       result = prop;
  return result;
}

// http://stackoverflow.com/questions/13303151/getting-fullscreen-mode-to-my-browser-using-jquery
function toggleFullScreen() {
  if (!document.fullscreenElement &&    // alternative standard method
      !document.mozFullScreenElement && !document.webkitFullscreenElement && !document.msFullscreenElement ) {  // current working methods
    if (document.documentElement.requestFullscreen) {
      document.documentElement.requestFullscreen();
    } else if (document.documentElement.msRequestFullscreen) {
      document.documentElement.msRequestFullscreen();
    } else if (document.documentElement.mozRequestFullScreen) {
      document.documentElement.mozRequestFullScreen();
    } else if (document.documentElement.webkitRequestFullscreen) {
      document.documentElement.webkitRequestFullscreen(Element.ALLOW_KEYBOARD_INPUT);
    }
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
  }
}

// http://www.davekoelle.com/files/alphanum.js
function alphanum(a, b) {
  function chunkify(t) {
    var tz = new Array();
    var x = 0, y = -1, n = 0, i, j;
    while (i = (j = t.charAt(x++)).charCodeAt(0)) {
      var m = (i == 46 || (i >=48 && i <= 57));
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
      } else return (aa[x] > bb[x]) ? 1 : -1;
    }
  }
  return aa.length - bb.length;
}

// https://github.com/os-js/OS.js/blob/30c7ccac294936dcd2a62c640e4759edb524841e/src/client/javascript/utils/fs.js
function dirname(f) {
  var pstr   = f.split(/^(.*)\:\/\/(.*)/).filter(function(n) { return n !== ''; });
  var args   = pstr.pop();
  var prot   = pstr.pop();
  var result = '';
  var tmp = args.split('/').filter(function(n) { return n !== ''; });
  if ( tmp.length ) {
    tmp.pop();
  }
  result = tmp.join('/');
  if ( !result.match(/^\//) ) {
    result = '/' + result;
  }
  if ( prot ) {
    result = prot + '://' + result;
  }
  return result;
}

function basename(str) {
  return new String(str).substring(str.lastIndexOf('/') + 1);
}

// http://www.mikedoesweb.com/2014/javascript-object-next-and-previous-keys/
var oFunctions = {};
oFunctions.keys = {};
//NEXT KEY
oFunctions.keys.next = function(o, id){
  var keys = Object.keys( o ).sort(alphanum);
  var idIndex = keys.indexOf( id );
  var nextIndex = idIndex += 1;
  if(nextIndex >= keys.length){
    //we're at the end, there is no next
    return;
  }
  var nextKey = keys[ nextIndex ]
  return nextKey;
};
//PREVIOUS KEY
oFunctions.keys.previous = function(o, id){
  var keys = Object.keys( o ).sort(alphanum);
  var idIndex = keys.indexOf( id ) + 1;
  var nextIndex = idIndex -= 2;
  var nextKey = keys[ nextIndex ]
  return nextKey;
};

