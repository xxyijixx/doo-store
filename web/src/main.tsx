import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import App from './App.tsx'

import { I18nextProvider } from 'react-i18next'
import i18n from './i18n.ts'



// import  { GlobalStoreProvider } from './GlobalStoreContext.tsx'

// const root = ReactDOM.createRoot(document.getElementById('root') as HTMLElement);

// root.render(
//   <GlobalStoreProvider>
//     <App />
//   </GlobalStoreProvider>
// );

// proxy document.body
const microAppBody = document.querySelector("micro-app-body");

if (microAppBody) {
  // Replace document.body with the micro-app-body element
  Object.defineProperty(document, "body", {
    get() {
      return microAppBody;
    },
  });

  // Example to show that it works
  console.log("document.body:", document.body);
} else {
  console.error("<micro-app-body> element not found!");
}

createRoot(document.getElementById("doo-store")!).render(
  <StrictMode>
    <I18nextProvider i18n={i18n}>
      <App />
    </I18nextProvider>
  </StrictMode>,
)
