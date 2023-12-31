
# justGame

A test e-commerce website designed as a game, allowing users to engage in buying and selling activities.

# Deployment

Follow these steps to deploy the project locally:

### 1. Clone Repository

```bash
git clone https://github.com/oOWinOo/justGame.git
```

### 2. Open Docker
###### Make sure Docker is installed on your machine.
####    Warning Port Usage
##### This project uses the following localhost ports:

- Port 8080: Web application
- Port 5432: PostgreSQL database
- Port 5050: pgAdmin for database administration
Ensure that no other processes are using these ports on your machine to avoid conflicts. If any of these ports are already in use, you may need to stop the conflicting processes or modify the Docker configuration.

### 3. Docker Compose Up

```bash
docker-compose up -d
```
### 4. run main.go

```bash
go run main.go
```


# Open Index Page
#### Visit localhost:8080/login to open the starting page "login page" and can select the register to create new account. 

![Login Screenshot](https://cdn.discordapp.com/attachments/1182874755565621350/1193843043296620554/image.png?ex=65ae2ffa&is=659bbafa&hm=e722ab771da7f8c5b5691d0b3c5645bd449efc35ef024bf70dd3e9ac9ba3e78c&)

![Register Screenshot](https://cdn.discordapp.com/attachments/1182874755565621350/1193843399929892954/image.png?ex=65ae304f&is=659bbb4f&hm=5383bbcaa03eead0f5852c8f4105031e7c9b3cb696f91170e3e390befab43e5b&)

#### After login will go to landing page that showing market that all user sell it on. 

![Market Screenshot](https://cdn.discordapp.com/attachments/1182874755565621350/1193843627793842176/image.png?ex=65ae3085&is=659bbb85&hm=b827117f45dcc7713813e78bb9d79b0101f791990eebbd52a542d1abd1cc8988&)

#### Spend the money to buy from market or Receive random products from system to increase login reward. 

![BuyMarket Screenshot](https://cdn.discordapp.com/attachments/1182874755565621350/1193843971676446841/image.png?ex=65ae30d7&is=659bbbd7&hm=9c1dd45d00299037a852b81bc4356440bc57bd09340f9f010c7b75962e63cedf&)


#### The products that you receive will be in inventory.

![Inventory Screenshot](https://cdn.discordapp.com/attachments/1182874755565621350/1193843715970711552/image.png?ex=65ae309a&is=659bbb9a&hm=87d9dd5607650aa29394c4985afbbdd3f13b12d55c1532f66d3eaaad114086b9&)

#### Sell products on the market that you can offer the cost.
![SellMarket Screenshot](https://cdn.discordapp.com/attachments/1182874755565621350/1193843783989727274/image.png?ex=65ae30aa&is=659bbbaa&hm=5477d8376d7d5e162c571930d9edfed5fe034483599c6e8e41fdcb1cdbbc6c4c&)

#### Your products are being sold in the market, but if no buyers come yet. They will be displayed on My Market and you can cancel to sell it anytime.
![MyMarket Screenshot](https://cdn.discordapp.com/attachments/1182874755565621350/1193843837576163388/image.png?ex=65ae30b7&is=659bbbb7&hm=ff22fd3890b61207089afcc4cdb3275bb1f1364e6a0372eef5590b182b28b50d&)

#### Another way to receive money from a product is to sell it to the system at the default price, and that product will be deleted.
![SellSystem Screenshot](https://cdn.discordapp.com/attachments/1182874755565621350/1193843926247944374/image.png?ex=65ae30cc&is=659bbbcc&hm=17896704efa0379e9d7d45360c3bde809ac9b6b9827b93d010a3d8ae9bc72251&)

#### At first, the inventory's capacity is 10, which means you can have only 10 products in your inventory. However, you can upgrade it with your money to expand it by 10 spaces.
![UpgradeInventory Screenshot](https://cdn.discordapp.com/attachments/1182874755565621350/1193844923523416094/image.png?ex=65ae31ba&is=659bbcba&hm=ee03a669513d774eec565468765b7aaab076fce13bf350f24becf57b3427d357&)

#### All your income and expenses will be displayed on the history.
![History Screenshot](https://cdn.discordapp.com/attachments/1182874755565621350/1193845773956300800/image.png?ex=65ae3285&is=659bbd85&hm=f22c8759890d5d50bf21505f30a6179e0c883d23adaffd20097454b1459910ca&)



# View data in PostgreSQL

### Access localhost:5050 for pgAdmin.

- Username: admin@email.com
- Password: admin

#### Set up the connection with the following details:

- Database Name: Anything you want
- Hostname: postgres
- Maintenance: postgres
- Username: postgres
- Password: postgres