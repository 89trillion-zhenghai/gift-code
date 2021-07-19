from locust import HttpUser, TaskSet, task,between


class QuickstartUser(HttpUser):
    wait_time = between(1, 2)

    @task
    def createAndGetGiftCode01(self):
       self.client.post('/createAndGetGiftCode',data={'userName':'admin01','description':'十周年活动奖励','giftType':'1','validity':'1m','availableTimes':'9999','giftDetail':'{"道具":"10","小兵":"20","金币":"10000"}'})

    @task
    def createAndGetGiftCode02(self):
       self.client.post('/createAndGetGiftCode',data={'userName':'admin01','description':'十周年活动奖励','giftType':'2','validity':'2m','availableTimes':'500000','giftDetail':'{"小兵":"20","钻石":"100"}'})

    @task
    def createAndGetGiftCode03(self):
       self.client.post('/createAndGetGiftCode',data={'userName':'admin01','description':'十周年活动奖励','giftType':'3','validity':'5m','availableTimes':'1','giftDetail':'{"道具":"20","小兵":"20","经验":"9000"}'})

    @task
    def getGiftDetail01(self):
       self.client.post("/getGiftDetail",data={'giftCode':'66O78E63'})

    @task
    def getGiftDetail02(self):
       self.client.post("/getGiftDetail",data={'giftCode':'PM5B651B'})

    @task
    def getGiftDetail03(self):
       self.client.post("/getGiftDetail",data={'giftCode':'22543S07'})

    @task
    def redeemGift01(self):
       self.client.post('/redeemGift',data={'giftCode':'66O78E63','uuid':'user01'})
    @task
    def redeemGift02(self):
       self.client.post('/redeemGift',data={'giftCode':'PM5B651B','uuid':'user01'})
    @task
    def redeemGift03(self):
       self.client.post('/redeemGift',data={'giftCode':'66O78E63','uuid':'user02'})
    @task
    def redeemGift04(self):
       self.client.post('/redeemGift',data={'giftCode':'22543S07','uuid':'user02'})
    @task
    def redeemGift05(self):
       self.client.post('/redeemGift',data={'giftCode':'PM5B651B','uuid':'user03'})
    @task
    def redeemGift06(self):
       self.client.post('/redeemGift',data={'giftCode':'66O78E63','uuid':'user03'})