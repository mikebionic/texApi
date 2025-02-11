
## Register process

1. Send register requests with **Credentials** and **RegisterMethod** in header
2. Validate email or phone number accordingly
3. If user not found or Verified != 1 -> 
   + save user with OTPKey VerifyTime, to track otp validation time
4. OTPKey validation, with email or phone included, if success -> update VerifyTime to track register process time
5. register user with POST request

## Google auth workflow

How to create OAuth and get required keys.
```text
go to console.cloud.google.com
go to App -> credentials -> create OAuth credential
Write type, name, add authorized URI's
copy data of ClientID and Secret to .env
```

Following text is legacy:
+ Client Registration: First, we need to register our application with Google. We should provide our application name, redirection URLs, etc. After registering, Google will provide us with a unique Client ID and Client Secret.
+ User Sign-In: When the user clicks the ‘Sign in with Google’ button on our website, they are redirected to Google’s authorization server.
+ Authorization Code: After the user is authenticated successfully, Google’s server will send a request to our application (to the URL we provided while registering our application with Google) along with an authorization code. Our application controller will receive the request from Google.
+ Access Token Request: Our application will then make a POST call to Google requesting an access token, attaching the authorization code received in the previous step, as well as the client ID and client secret (received during client registration).
+ User Details Retrieval: Once we obtain the access token, we can request Google to provide the user’s details by authorizing the access token.
+ Profile Creation: After we get the user’s details, such as email, username, and profile picture from Google, we can create a profile for the user in our application without requiring them to manually provide their details.

## Recovery password ForgotPassword workflow

+ /password/forgot/ request with Credentials CredType in header
+ /validate-otp/
+ /password/update/

----
Логистическая кампания перевозчик и как лог кампания делают поиск груза и предоставляет услугу водителя.

А отправитель создает груз и ищет водителей.
Они оба делают GET POST в Offer
просто в одном случае они в качестве
offer_role: sender / carrier
и указывать свой cargo_id, vehicle_id.
user_id, company_id также берется с пользователя, каким он бы ни был.

потом водители могут откликнуться и выбрать Оферту которую хотят перевезти, или логистическая кампания уже на момент выбора груза, предлагает своего водителя и автомобиль.

offer_requests, которые как ты говорил нужно смотреть заявки и сделать кого то победителем, чтобы их данные автоматом перетекли в нашу оферту.