function insert(){
    var chV = document.getElementById('todoV').value
    if (chV === "") {
    window.alert("請勿輸入空值")
    }else{
      const uri = `/insert`;
      const initDetails = {
          method: 'POST',
          headers: {
              "Content-Type": "application/json; charset=utf-8"
          },
          body: JSON.stringify({"title":chV}),
          mode: "cors"
      }
          fetch( uri, initDetails )
              .then(function (response) {
                return response.json();
              })
              .then(function (data) {
              // window.alert(JSON.stringify(data.title)+' 新增成功',
              location.reload()
             })
              .catch( err =>
              {
                  console.log( 'Fetch Error :-S', err );
              } );
    }
    }
    function change(e){
      var chV = prompt('請輸入修改值');
      var chId = e.value
      var cid = parseInt(chId,10)
      const uri = `/put/${chId}`;
      const initDetails = {
          method: 'PUT',
          headers: {
              "Content-Type": "application/json; charset=utf-8"
          },
          body: JSON.stringify({"id":cid,"title":chV}),
          mode: "cors"
      }
          fetch( uri, initDetails )
              .then(function (response) {
                return response.json();
              })
              .then(function (data) {
              appendData(data);
              // window.alert(JSON.stringify(data.title)+'修改成功')
              location.reload()
             })
              .catch( err =>
              {
                  console.log( 'Fetch Error :-S', err );
              } );
    }
    
    function del(e){
      var chId = e.value
      var cid = parseInt(chId,10)
      const uri = `/del/${chId}`;
      const initDetails = {
          method: 'DELETE',
          headers: {
              "Content-Type": "application/json; charset=utf-8"
          },
          body: JSON.stringify({"id":cid}),
          mode: "cors"
      }
          fetch( uri, initDetails )
              .then(function (response) {
                return response.json();
              })
              .then(function (data) {
              appendData(data);
              // window.alert(JSON.stringify(data))
              location.reload()
              
             })
              .catch( err =>
              {
                  console.log( 'Fetch Error :-S', err );
              } );
    }
    function finish(e){
      if (e.value =="false"){
       var a = document.getElementsByClassName('finish').value = true
      }else{
       var a = document.getElementsByClassName('finish').value = false
      }
      var chId = e.id
      var cid = parseInt(chId,10)
      const uri = `/put/${chId}`;
      const initDetails = {
          method: 'PUT',
          headers: {
              "Content-Type": "application/json; charset=utf-8"
          },
          body: JSON.stringify({"id":cid,"status":a}),
          mode: "cors"
      }
          fetch( uri, initDetails )
              .then(function (response) {
                return response.json();
              })
              .then(function (data) {
                location.reload()
              appendData(data);
             })
              .catch( err =>
              {
                  console.log( 'Fetch Error :-S', err );
              } );
              function appendData(data) {
                // var a = JSON.stringify(data.status)
                if (data.status == true){
                  alert(data.title+"已完成")
                }else{
                  alert(data.title+"未完成")
    
                }
              }
    }