jQuery(function($)  {
sock.onopen  =  function()  {
	if( debug ) console.log("connected  to  "  +  wsuri);
}
sock.onclose  =  function(e)  {
	if( debug ) console.log("connection  closed  ("  +  e.code  +  ")");
}
sock.onmessage  =  function(e)  {
	var  data  =  JSON.parse(e.data),
		img = null;
	switch(data.Type)  {
		case  0://离开
			user_leave(data.Uid,data.Username);
		break;
		case  1://加入
			user_join(data.Uid,data.Username);
			break;
		case  2://发送信息
			img = user_img(data.Uid);
			send_msg(data.Username,img,data.Msg,true);
			new_msg();
		break;

	}
	if( debug ) console.log("message  received:  "  +  e.data);
}

var user_template = null;

$('#send').click(function(){
	var  msg  =  $('#msg-box').val();
	if(  msg  !=  ''){
		sock.send(msg);
	}
});
  
$('.chat-message  input').keypress(function(e){
	if(e.which  ==  13)  {
		var  msg  =  $('#msg-box').val();
		if(  msg  !=  ''){
			sock.send(msg);
		}
	}
});

var  i  =  0;
	msg = null;
function  send_msg(name,img,msg,clear)  {
	i  =  i  +  1;
	var    inner  =  $('#chat-messages-inner');
	var  time  =  new  Date();
	var  hours  =  time.getHours();
	var  minutes  =  time.getMinutes();
	if(hours  <  10)  hours  =  '0'  +  hours;
	if(minutes  <  10)  minutes  =  '0'  +  minutes;
	var  id  =  'msg-'+i;
	var  idname  =  name.replace('  ','-').toLowerCase();
	inner.append('<p  id="'+id+'"  class="user-'+idname+'">'
	                                +'<span  class="msg-block"><img  src="'+img+'"  alt=""  /><strong>'+name+'</strong>  <span  class="time">-  '+hours+':'+minutes+'</span>'
	                                +'<span  class="msg">'+msg+'</span></span></p>');
	$('#'+id).hide().fadeIn(800);
	if(clear)  {
		$('.chat-message  input').val('').focus();
	}
	$('#chat-messages').animate({  scrollTop:  inner.height()  },1000);
}
function user_join(userid,name) {
	user_template  =  '<li  id="user_'+userid+'"  class="online"><a  href="#"><img  alt=""  src="'+user_img(userid)+'"  /><span>'+name+'</span></a></li>';
	$('#user_list').append(user_template);
}
function user_leave(userid,name)  {
	i  =  i  +  1;
	$('#user_'+userid).addClass('offline').delay(1000).slideUp(800,function(){
		$(this).remove();
	});
	var    inner  =  $('#chat-messages-inner');
	var  id  =  'msg-'+i;
	inner.append('<p  class="offline"  id="'+id+'"><span>用户  '+name+'  离开聊天室</span></p>');
	$('#'+id).hide().fadeIn(800);
}
function user_img(userid) {
	userid = userid % 5 + 1;
	return '/static/img/av'+userid+'.jpg';
}
});