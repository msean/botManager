/*! 
 Build based on gin-vue-admin 
 Time : 1766128845000 */
import{_ as e}from"./087AC4D233B64EB0_plugin-vue_export-helper.BCo6x5W8.js";import{J as s,c as a,o as i,d as o}from"./087AC4D233B64EB0index.D6EEDy3E.js";const l={class:"media"},p=["src"],t={key:1,class:"sticker-placeholder"},r=e({__name:"ChatMedia",props:{msg:Object},setup(e){const r=e,c=s(()=>"photo"===r.msg.fileType);s(()=>"video"===r.msg.fileType),s(()=>"voice"===r.msg.fileType),s(()=>"audio"===r.msg.fileType);const m=s(()=>"sticker"===r.msg.fileType);return(s,r)=>(i(),a("div",l,[c.value?(i(),a("img",{key:0,src:e.msg.fileUrl,class:"photo"},null,8,p)):m.value?(i(),a("div",t," <表情包> ")):o("",!0)]))}},[["__scopeId","data-v-3a46a08b"]]);export{r as default};
