from rest_framework import serializers

from products.models import Product


class UserSerializer(serializers.HyperlinkedModelSerializer):
    class Meta:
        model = Product
        fields = ['id', 'name', 'link', 'price']
