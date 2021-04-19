from rest_framework import serializers

from products.models import Products


class UserSerializer(serializers.HyperlinkedModelSerializer):
    class Meta:
        model = Products
        fields = ['id', 'name', 'price']
